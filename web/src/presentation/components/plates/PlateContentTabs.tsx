"use client"

import { Fragment, useState, useEffect, useMemo } from "react"
import type { ComponentPropsWithoutRef } from "react"
import { ChevronDown, ChevronRight, File, Folder, Loader2 } from "lucide-react"
import type { RepoTreeEntry } from "@/src/data/repositories/githubClient"
import { resolveRepoMarkdownHref } from "@/src/presentation/utils/readmeLinks"
import { AuthService } from "@/src/domain/services/AuthService"
import { toast } from "sonner"

interface Props {
  readme: string | null
  license: string | null
  tree: RepoTreeEntry[] | null
  schemaFields?: Array<{ key: string; type: string; required?: boolean; defaultValue?: string; values?: string[] }>
  modules?: Array<{ name: string; enabled?: boolean }>
  hasGenerate?: boolean
  slug?: string
  repoUrl?: string | null
  branch?: string | null
}

type Tab = "readme" | "license" | "files" | "schema"

interface FileTreeNode {
  name: string
  path: string
  type: "blob" | "tree"
  children: FileTreeNode[]
}

interface BuilderNode {
  name: string
  path: string
  type: "blob" | "tree"
  children: Map<string, BuilderNode>
}

function buildFileTree(entries: Array<{ path?: string; type: "blob" | "tree" }>): FileTreeNode[] {
  const root: BuilderNode = {
    name: "",
    path: "",
    type: "tree",
    children: new Map(),
  }

  const sorted = entries
    .filter((entry): entry is { path: string; type: "blob" | "tree" } => typeof entry.path === "string" && entry.path.trim() !== "")
    .sort((a, b) => a.path.localeCompare(b.path))

  for (const entry of sorted) {
    const parts = entry.path.split("/").filter(Boolean)
    if (parts.length === 0) continue

    let current = root
    let currentPath = ""

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i]
      currentPath = currentPath ? `${currentPath}/${part}` : part
      const isLeaf = i === parts.length - 1

      const existing = current.children.get(part)
      if (existing) {
        if (!isLeaf) existing.type = "tree"
        if (isLeaf) existing.type = entry.type
        current = existing
        continue
      }

      const created: BuilderNode = {
        name: part,
        path: currentPath,
        type: isLeaf ? entry.type : "tree",
        children: new Map(),
      }

      current.children.set(part, created)
      current = created
    }
  }

  const toArray = (node: BuilderNode): FileTreeNode[] => {
    const items = Array.from(node.children.values())
      .sort((a, b) => {
        if (a.type !== b.type) return a.type === "tree" ? -1 : 1
        return a.name.localeCompare(b.name)
      })
      .map((child) => ({
        name: child.name,
        path: child.path,
        type: child.type,
        children: toArray(child),
      }))

    return items
  }

  return toArray(root)
}

function collectFolderPaths(nodes: FileTreeNode[]): string[] {
  const folderPaths: string[] = []

  const walk = (items: FileTreeNode[]) => {
    for (const item of items) {
      if (item.type === "tree") {
        folderPaths.push(item.path)
        if (item.children.length > 0) {
          walk(item.children)
        }
      }
    }
  }

  walk(nodes)
  return folderPaths
}

function normalizeType(raw: string): string {
  return raw.trim().toLowerCase()
}

function coerceInputValue(type: string, rawValue: string): unknown {
  const normalized = normalizeType(type)

  if (normalized === "bool" || normalized === "boolean") {
    return rawValue.toLowerCase() === "true"
  }

  if (normalized === "int" || normalized === "integer") {
    const n = Number.parseInt(rawValue, 10)
    return Number.isNaN(n) ? rawValue : n
  }

  if (normalized === "number" || normalized === "float" || normalized === "double") {
    const n = Number(rawValue)
    return Number.isNaN(n) ? rawValue : n
  }

  return rawValue
}

function TreeView({
  nodes,
  expandedPaths,
  onToggle,
  depth = 0,
}: {
  nodes: FileTreeNode[]
  expandedPaths: Set<string>
  onToggle: (path: string) => void
  depth?: number
}) {
  return (
    <ul className="space-y-1">
      {nodes.map((node) => (
        <li key={node.path}>
          <div
            className="flex items-center gap-2 py-0.5 text-sm"
            style={{ paddingLeft: `${depth * 16}px` }}
          >
            {node.type === "tree" ? (
              <button
                type="button"
                onClick={() => onToggle(node.path)}
                className="inline-flex items-center gap-1 text-left"
                aria-label={`${expandedPaths.has(node.path) ? "Collapse" : "Expand"} ${node.name}`}
              >
                {expandedPaths.has(node.path) ? (
                  <ChevronDown className="h-3.5 w-3.5 text-muted-foreground" />
                ) : (
                  <ChevronRight className="h-3.5 w-3.5 text-muted-foreground" />
                )}
                <Folder className="h-3.5 w-3.5 text-muted-foreground" />
                <span className="font-mono text-xs text-foreground">{node.name}</span>
              </button>
            ) : (
              <>
                <span className="inline-block w-3.5" />
                <File className="h-3.5 w-3.5 text-muted-foreground" />
                <span className="font-mono text-xs text-foreground">{node.name}</span>
              </>
            )}
          </div>
          {node.type === "tree" && node.children.length > 0 && expandedPaths.has(node.path) ? (
            <TreeView
              nodes={node.children}
              expandedPaths={expandedPaths}
              onToggle={onToggle}
              depth={depth + 1}
            />
          ) : null}
        </li>
      ))}
    </ul>
  )
}

export function PlateContentTabs({
  readme,
  license,
  tree,
  schemaFields = [],
  modules = [],
  hasGenerate: hasGenerateProp = false,
  slug,
  repoUrl,
  branch,
}: Props) {
  const [active, setActive] = useState<Tab>("readme")
  const [MarkdownComponent, setMarkdownComponent] = useState<React.ComponentType<Record<string, unknown>> | null>(null)
  const [plugins, setPlugins] = useState<unknown[]>([])
  const treeNodes = useMemo(() => (tree ? buildFileTree(tree) : []), [tree])
  const folderPaths = useMemo(() => collectFolderPaths(treeNodes), [treeNodes])
  const [expandedPaths, setExpandedPaths] = useState<Set<string>>(new Set())
  const hasGenerate = hasGenerateProp || schemaFields.length > 0
  const [downloadValues, setDownloadValues] = useState<Record<string, string>>({})
  const [moduleValues, setModuleValues] = useState<Record<string, boolean>>({})
  const [isGenerating, setIsGenerating] = useState(false)

  useEffect(() => {
    const topLevelFolders = treeNodes
      .filter((node) => node.type === "tree")
      .map((node) => node.path)
    setExpandedPaths(new Set(topLevelFolders))
  }, [treeNodes])

  const toggleFolder = (path: string) => {
    setExpandedPaths((prev) => {
      const next = new Set(prev)
      if (next.has(path)) {
        next.delete(path)
      } else {
        next.add(path)
      }
      return next
    })
  }

  const expandAll = () => setExpandedPaths(new Set(folderPaths))
  const collapseAll = () => setExpandedPaths(new Set())

  useEffect(() => {
    Promise.all([
      import("react-markdown"),
      import("remark-gfm"),
      import("remark-breaks"),
    ]).then(([md, gfm, breaks]) => {
      setMarkdownComponent(() => md.default)
      setPlugins([gfm.default, breaks.default])
    })
  }, [])

  useEffect(() => {
    const next: Record<string, string> = {}
    for (const field of schemaFields) {
      next[field.key] = field.defaultValue ?? ""
    }
    setDownloadValues(next)
  }, [schemaFields])

  useEffect(() => {
    const next: Record<string, boolean> = {}
    for (const mod of modules) {
      next[mod.name] = Boolean(mod.enabled)
    }
    setModuleValues(next)
  }, [modules])

  useEffect(() => {
    const syncFromHash = () => {
      if (window.location.hash === "#files" && tree && tree.length > 0) {
        setActive("files")
        return
      }
      if (window.location.hash === "#schema" && hasGenerate) {
        setActive("schema")
        return
      }
      if (window.location.hash === "#license" && license) {
        setActive("license")
        return
      }
      if (window.location.hash === "#readme" && readme) {
        setActive("readme")
        return
      }

      if (readme) {
        setActive("readme")
        return
      }
      if (license) {
        setActive("license")
        return
      }
      if (hasGenerate) {
        setActive("schema")
        return
      }
      if (tree && tree.length > 0) {
        setActive("files")
      }
    }

    syncFromHash()
    window.addEventListener("hashchange", syncFromHash)
    return () => window.removeEventListener("hashchange", syncFromHash)
  }, [readme, license, tree, hasGenerate])

  const content = active === "readme" ? readme : active === "license" ? license : null

  const handleGenerateDownload = async () => {
    const missingRequired = schemaFields
      .filter((field) => field.required)
      .filter((field) => !(downloadValues[field.key] ?? "").trim())

    if (missingRequired.length > 0) {
      toast.error(`Please fill required fields: ${missingRequired.map((f) => f.key).join(", ")}`)
      return
    }

    const values: Record<string, unknown> = {}
    for (const field of schemaFields) {
      const raw = (downloadValues[field.key] ?? "").trim()
      if (!raw) continue
      values[field.key] = coerceInputValue(field.type, raw)
    }
    for (const mod of modules) {
      values[`modules.${mod.name}.enabled`] = Boolean(moduleValues[mod.name])
    }

    try {
      setIsGenerating(true)
      const token = AuthService.getToken()
      const slugFromPath = slug || window.location.pathname.split("/").filter(Boolean).pop() || "plate"
      const res = await fetch(`/api/generate/${encodeURIComponent(slugFromPath)}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({ values }),
      })

      if (!res.ok) {
        let message = "Failed to generate project"
        try {
          const data = await res.json() as { error?: string }
          if (data.error) message = data.error
        } catch {}
        throw new Error(message)
      }

      const blob = await res.blob()
      const objectUrl = window.URL.createObjectURL(blob)
      const anchor = document.createElement("a")
      anchor.href = objectUrl
      anchor.download = `${slugFromPath}.zip`
      document.body.appendChild(anchor)
      anchor.click()
      document.body.removeChild(anchor)
      window.URL.revokeObjectURL(objectUrl)
      toast.success("Project generated and download started")
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to generate project"
      toast.error(message)
    } finally {
      setIsGenerating(false)
    }
  }

  const branchResolved = (branch && branch.trim()) || "main"
  const githubRepo = repoUrl?.trim() ? repoUrl : undefined
  const sourceFileInRepo = active === "license" ? "LICENSE" : "README.md"

  const markdownComponents = useMemo(
    () => ({
      pre: ({ children, ...props }: ComponentPropsWithoutRef<"pre">) => (
        <pre
          {...props}
          className="my-4 overflow-x-auto border !border-border !bg-muted dark:!bg-card !p-6 text-xs leading-relaxed !text-foreground"
        >
          {children}
        </pre>
      ),
      code: ({ inline, className, children, ...props }: ComponentPropsWithoutRef<"code"> & { inline?: boolean }) => {
        if (inline) {
          return (
            <code
              {...props}
              className="rounded-none border border-border bg-muted/80 px-2 py-1 font-mono text-xs text-foreground dark:bg-card"
            >
              {children}
            </code>
          )
        }

        return (
          <code
            {...props}
            className={className ? `${className} !text-foreground font-mono` : "!text-foreground font-mono"}
          >
            {children}
          </code>
        )
      },
      a: ({ href, children, ...props }: ComponentPropsWithoutRef<"a">) => {
        const resolved = resolveRepoMarkdownHref(href, githubRepo, branchResolved, sourceFileInRepo)
        const openNewTab = /^https?:\/\//i.test(resolved)
        return (
          <a
            {...props}
            href={resolved}
            {...(openNewTab ? { target: "_blank", rel: "noopener noreferrer" } : {})}
          >
            {children}
          </a>
        )
      },
      img: ({ src, alt, ...props }: ComponentPropsWithoutRef<"img">) => {
        const srcStr = typeof src === "string" ? src : undefined
        const resolved = srcStr
          ? resolveRepoMarkdownHref(srcStr, githubRepo, branchResolved, sourceFileInRepo)
          : srcStr
        return <img {...props} src={resolved} alt={alt ?? ""} />
      },
    }),
    [githubRepo, branchResolved, sourceFileInRepo],
  )

  return (
    <div>
      <div>
        {active === "schema" ? (
          hasGenerate ? (
            <div className="space-y-5 border border-border bg-card p-6">
              {schemaFields.length > 0 ? (
                <div>
                  <p className="mb-2 text-sm font-semibold text-foreground">Schema fields</p>
                  <div className="overflow-auto border border-border bg-background">
                    <table className="w-full min-w-[640px] text-xs">
                      <thead className="bg-card text-muted-foreground">
                        <tr>
                          <th className="px-3 py-2 text-left font-medium">Field</th>
                          <th className="px-3 py-2 text-left font-medium">Type</th>
                          <th className="px-3 py-2 text-left font-medium">Required</th>
                          <th className="px-3 py-2 text-left font-medium">Default</th>
                        </tr>
                      </thead>
                      <tbody>
                        {schemaFields.map((field, index) => {
                          const isEven = index % 2 === 0
                          const tone = isEven ? "bg-background" : "bg-muted/60 dark:bg-card"
                          const hasValues = Boolean(field.values && field.values.length > 0)

                          return (
                            <Fragment key={field.key}>
                              <tr key={`${field.key}-main`} className={`border-t border-border/70 ${tone}`}>
                                <td className="px-3 pb-0.5 pt-2 align-top">
                                  <div className="font-mono text-[13px] text-foreground">{field.key}</div>
                                </td>
                                <td className="px-3 pb-0.5 pt-2 align-top text-foreground">{field.type}</td>
                                <td className="px-3 pb-0.5 pt-2 align-top text-foreground">{field.required ? "yes" : "no"}</td>
                                <td className="px-3 pb-0.5 pt-2 align-top text-muted-foreground">{field.defaultValue || "-"}</td>
                              </tr>
                              {hasValues ? (
                                <tr key={`${field.key}-meta`} className={tone}>
                                  <td colSpan={4} className="px-3 pb-2 pt-0.5 text-muted-foreground">
                                    <div className="flex flex-wrap gap-2">
                                      {field.values!.map((value) => (
                                        <span
                                          key={`${field.key}-${value}`}
                                          className="border border-border/70 bg-background px-1.5 py-0.5 font-mono text-[10px] text-foreground"
                                        >
                                          {value}
                                        </span>
                                      ))}
                                    </div>
                                  </td>
                                </tr>
                              ) : null}
                            </Fragment>
                          )
                        })}
                      </tbody>
                    </table>
                  </div>
                </div>
              ) : null}

              {modules.length > 0 ? (
                <div>
                  <p className="mb-2 text-sm font-semibold text-foreground">Modules</p>
                  <div className="flex flex-wrap gap-2">
                    {modules.map((mod) => (
                      <span key={mod.name} className="border border-border bg-background px-2 py-1 text-xs text-foreground">
                        <span className="font-mono">{mod.name}</span>
                        <span className="ml-1 text-muted-foreground">default {mod.enabled ? "on" : "off"}</span>
                      </span>
                    ))}
                  </div>
                </div>
              ) : null}

              <div className="space-y-3 border border-border bg-background p-3">
                <p className="text-sm font-semibold text-foreground">Build and download project</p>
                <div className="grid gap-2 sm:grid-cols-2">
                  {schemaFields.map((field) => {
                    const normalized = normalizeType(field.type)
                    const hasEnum = Array.isArray(field.values) && field.values.length > 0

                    return (
                      <label key={`input-${field.key}`} className="grid gap-1">
                        <span className="text-[11px] text-muted-foreground">
                          {field.key}
                          {field.required ? " *" : ""}
                        </span>
                        {hasEnum ? (
                          <select
                            className="h-8 border border-border bg-background px-2 text-xs"
                            value={downloadValues[field.key] ?? ""}
                            onChange={(e) => setDownloadValues((prev) => ({ ...prev, [field.key]: e.target.value }))}
                          >
                            <option value="">Select value</option>
                            {field.values!.map((value) => (
                              <option key={`${field.key}-${value}`} value={value}>
                                {value}
                              </option>
                            ))}
                          </select>
                        ) : normalized === "bool" || normalized === "boolean" ? (
                          <select
                            className="h-8 border border-border bg-background px-2 text-xs"
                            value={downloadValues[field.key] ?? ""}
                            onChange={(e) => setDownloadValues((prev) => ({ ...prev, [field.key]: e.target.value }))}
                          >
                            <option value="">Select value</option>
                            <option value="true">true</option>
                            <option value="false">false</option>
                          </select>
                        ) : (
                          <input
                            className="h-8 border border-border bg-background px-2 text-xs"
                            type={(normalized === "int" || normalized === "integer" || normalized === "number" || normalized === "float" || normalized === "double") ? "number" : "text"}
                            step={normalized === "int" || normalized === "integer" ? "1" : "any"}
                            value={downloadValues[field.key] ?? ""}
                            placeholder={field.defaultValue ?? ""}
                            onChange={(e) => setDownloadValues((prev) => ({ ...prev, [field.key]: e.target.value }))}
                          />
                        )}
                      </label>
                    )
                  })}
                </div>
                {modules.length > 0 ? (
                  <div className="space-y-2 border-t border-border pt-3">
                    <p className="text-xs font-semibold text-foreground">Modules</p>
                    <div className="grid gap-2 sm:grid-cols-2">
                      {modules.map((mod) => (
                        <label
                          key={`toggle-${mod.name}`}
                          className="flex h-8 items-center justify-between border border-border bg-background px-2 text-xs"
                        >
                          <span className="font-mono text-foreground">{mod.name}</span>
                          <input
                            type="checkbox"
                            checked={Boolean(moduleValues[mod.name])}
                            onChange={(e) =>
                              setModuleValues((prev) => ({
                                ...prev,
                                [mod.name]: e.target.checked,
                              }))
                            }
                          />
                        </label>
                      ))}
                    </div>
                  </div>
                ) : null}
                <button
                  type="button"
                  onClick={handleGenerateDownload}
                  disabled={isGenerating}
                  className="inline-flex h-8 items-center gap-2 border border-border bg-background px-3 text-xs font-semibold text-foreground transition-colors hover:bg-muted disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {isGenerating ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : null}
                  {isGenerating ? "Generating..." : "Generate and download ZIP"}
                </button>
              </div>
            </div>
          ) : (
            <div className="p-6">
              <p className="py-4 text-center text-sm text-muted-foreground">No generation schema found</p>
            </div>
          )
        ) : active === "files" ? (
          treeNodes.length > 0 ? (
            <div className="border border-border bg-card p-6">
              <div className="mb-4 flex items-center justify-end gap-2 border-b border-border pb-3">
                <button
                  type="button"
                  onClick={expandAll}
                  className="h-8 border border-border px-3 text-xs font-medium text-foreground transition-colors hover:bg-muted"
                >
                  Expand all
                </button>
                <button
                  type="button"
                  onClick={collapseAll}
                  className="h-8 border border-border px-3 text-xs font-medium text-foreground transition-colors hover:bg-muted"
                >
                  Collapse all
                </button>
              </div>
              <TreeView
                nodes={treeNodes}
                expandedPaths={expandedPaths}
                onToggle={toggleFolder}
              />
            </div>
          ) : (
            <div className="p-6">
              <p className="py-4 text-center text-sm text-muted-foreground">No file tree found</p>
            </div>
          )
        ) : content ? (
          MarkdownComponent ? (
            <div className="prose prose-neutral dark:prose-invert max-w-none p-6
              prose-headings:font-semibold
              prose-headings:border-b prose-headings:border-border prose-headings:pb-2 prose-headings:mb-4
              prose-h1:text-2xl prose-h2:text-xl prose-h3:text-lg
              prose-p:text-sm prose-p:leading-relaxed prose-p:text-foreground
              prose-a:text-blue-500 prose-a:no-underline hover:prose-a:underline
              prose-code:before:content-none prose-code:after:content-none
              prose-pre:bg-transparent prose-pre:border-0 prose-pre:p-0 prose-pre:shadow-none
              prose-blockquote:border-l-4 prose-blockquote:border-border prose-blockquote:text-muted-foreground prose-blockquote:not-italic
              prose-table:text-sm prose-th:text-left prose-th:font-semibold prose-th:border prose-th:border-border prose-th:px-3 prose-th:py-2 prose-th:bg-muted
              prose-td:border prose-td:border-border prose-td:px-3 prose-td:py-2
              prose-img:border prose-img:border-border
              prose-hr:border-border
              prose-li:text-sm prose-li:marker:text-muted-foreground
              prose-strong:text-foreground prose-strong:font-semibold
            ">
              <MarkdownComponent remarkPlugins={plugins} components={markdownComponents}>
                {content}
              </MarkdownComponent>
            </div>
          ) : (
            <div className="p-6">
              <p className="text-xs text-muted-foreground">Loading…</p>
            </div>
          )
        ) : (
          <div className="p-6">
            <p className="text-sm text-muted-foreground text-center py-4">
              {active === "readme" ? "No README found" : "No LICENSE found"}
            </p>
          </div>
        )}
      </div>
    </div>
  )
}