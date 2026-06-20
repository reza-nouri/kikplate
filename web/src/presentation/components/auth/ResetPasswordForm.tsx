"use client"

import { useMemo, useState } from "react"
import Link from "next/link"
import { useRouter, useSearchParams } from "next/navigation"
import { useResetPassword } from "@/src/presentation/hooks/useAuth"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

export function ResetPasswordForm() {
  const router = useRouter()
  const params = useSearchParams()
  const resetPassword = useResetPassword()
  const token = useMemo(() => params.get("token")?.trim() ?? "", [params])
  const [newPassword, setNewPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")

  const passwordsMatch = newPassword === confirmPassword

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!token) {
      toast.error("Password reset token is missing")
      return
    }

    if (!passwordsMatch) {
      toast.error("Password confirmation does not match")
      return
    }

    try {
      const result = await resetPassword.mutateAsync({ token, newPassword })
      toast.success(result.message)
      router.push("/login")
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Unable to reset password")
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="new-password">New password</Label>
        <Input
          id="new-password"
          type="password"
          placeholder="••••••••"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          required
          autoComplete="new-password"
          minLength={8}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="confirm-password">Confirm new password</Label>
        <Input
          id="confirm-password"
          type="password"
          placeholder="••••••••"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
          autoComplete="new-password"
          minLength={8}
        />
      </div>

      <Button type="submit" className="w-full" disabled={resetPassword.isPending || !passwordsMatch}>
        {resetPassword.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Reset password
      </Button>

      <p className="text-center text-sm text-muted-foreground">
        Back to{" "}
        <Link href="/login" className="text-foreground underline underline-offset-4 hover:text-primary">
          Sign in
        </Link>
      </p>
    </form>
  )
}
