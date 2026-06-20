import { ForgotPasswordForm } from "@/src/presentation/components/auth/ForgotPasswordForm"

export default function ForgotPasswordPage() {
  return (
    <div className="container mx-auto flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-sm border border-border bg-card p-6">
        <h1 className="mb-3 text-2xl font-bold">Forgot password</h1>
        <p className="mb-6 text-sm text-muted-foreground">Enter your email and we will send a reset link.</p>
        <ForgotPasswordForm />
      </div>
    </div>
  )
}
