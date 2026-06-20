"use client"

import { useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { useRequestPasswordReset } from "@/src/presentation/hooks/useAuth"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

export function ForgotPasswordForm() {
  const router = useRouter()
  const requestReset = useRequestPasswordReset()
  const [email, setEmail] = useState("")

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      const result = await requestReset.mutateAsync(email.trim().toLowerCase())
      toast.success(result.message)
      router.push("/login")
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Unable to request password reset")
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="email">Email</Label>
        <Input
          id="email"
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          autoComplete="email"
        />
      </div>

      <Button type="submit" className="w-full" disabled={requestReset.isPending}>
        {requestReset.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Send reset link
      </Button>

      <p className="text-center text-sm text-muted-foreground">
        Remembered your password?{" "}
        <Link href="/login" className="text-foreground underline underline-offset-4 hover:text-primary">
          Sign in
        </Link>
      </p>
    </form>
  )
}
