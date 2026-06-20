"use client"

import Link from "next/link"
import { useMounted } from "@/src/presentation/hooks/useMounted"
import { useMe } from "@/src/presentation/hooks/useAuth"
import { LoadingSpinner } from "@/src/presentation/components/common/LoadingSpinner"
import { SubmitRepositoryForm } from "@/src/presentation/components/submit/SubmitRepositoryForm"

export default function SubmitPage() {
  const mounted = useMounted()
  const { data: me, isLoading } = useMe()

  if (!mounted || isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <LoadingSpinner />
      </div>
    )
  }

  if (!me) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center space-y-4">
          <p className="text-muted-foreground">
            You need to sign in to submit a plate.
          </p>
          <Link
            href="/login"
            className="text-sm font-medium underline underline-offset-4"
          >
            Sign in →
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="border-b border-border bg-muted/10">
        <div className="container mx-auto max-w-7xl px-4 py-8">
          <h1 className="text-2xl font-bold text-foreground">Submit a plate</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Share a reusable template others can generate into ready-to-use projects.
          </p>
        </div>
      </div>

      <div className="container mx-auto max-w-7xl px-4 py-10">
        <SubmitRepositoryForm />
      </div>
    </div>
  )
}