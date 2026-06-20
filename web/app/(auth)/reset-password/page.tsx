import { Suspense } from "react"
import { ResetPasswordForm } from "@/src/presentation/components/auth/ResetPasswordForm"

export default function ResetPasswordPage() {
  return (
    <Suspense
      fallback={
        <div className="container mx-auto flex min-h-screen items-center justify-center px-4">
          <div className="w-full max-w-sm border border-border bg-card p-6">
            <h1 className="mb-3 text-2xl font-bold">Reset password</h1>
            <p className="mb-6 text-sm text-muted-foreground">Loading...</p>
          </div>
        </div>
      }
    >
      <div className="container mx-auto flex min-h-screen items-center justify-center px-4">
        <div className="w-full max-w-sm border border-border bg-card p-6">
          <h1 className="mb-3 text-2xl font-bold">Reset password</h1>
          <p className="mb-6 text-sm text-muted-foreground">Set a new password for your account.</p>
          <ResetPasswordForm />
        </div>
      </div>
    </Suspense>
  )
}
