"use client";

import { useState, FormEvent } from "react";

type ViewMode = "login" | "forgot-password";

export default function LoginPage() {
  const [role, setRole] = useState<"student" | "admin">("student");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [loggedUser, setLoggedUser] = useState("");

  // ── Forgot Password state ──
  const [view, setView] = useState<ViewMode>("login");
  const [resetEmail, setResetEmail] = useState("");
  const [resetNewPassword, setResetNewPassword] = useState("");
  const [resetShowPassword, setResetShowPassword] = useState(false);
  const [resetError, setResetError] = useState("");
  const [resetSuccess, setResetSuccess] = useState("");
  const [resetLoading, setResetLoading] = useState(false);

  // ── Login Submit ──
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ role, username, password }),
      });

      const data = await res.json();

      if (!res.ok) {
        setError(data.message || "Login gagal. Silakan coba lagi.");
        setLoading(false);
        return;
      }

      // Simpan Bearer Token ke localStorage
      localStorage.setItem("bearer_token", data.bearer_token);
      localStorage.setItem("user_role", data.role);
      localStorage.setItem("user_email", data.username);

      // ✅ BUKTI: Bearer Token berhasil diterima dari backend
      console.log("✅ Bearer Token diterima:", data.bearer_token);

      setLoggedUser(data.username);
      setSuccess(true);
      setLoading(false);
    } catch {
      setError("Tidak dapat terhubung ke server. Pastikan backend berjalan di port 8080.");
      setLoading(false);
    }
  };

  // ── Reset Password Submit ──
  const handleResetPassword = async (e: FormEvent) => {
    e.preventDefault();
    setResetError("");
    setResetSuccess("");
    setResetLoading(true);

    try {
      const res = await fetch("http://localhost:8080/api/reset-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: resetEmail, new_password: resetNewPassword }),
      });

      const data = await res.json();

      if (!res.ok) {
        setResetError(data.message || "Gagal mengubah password.");
        setResetLoading(false);
        return;
      }

      setResetSuccess(data.message || "Password berhasil diubah!");
      setResetLoading(false);
    } catch {
      setResetError("Tidak dapat terhubung ke server. Pastikan backend berjalan di port 8080.");
      setResetLoading(false);
    }
  };

  // ── Switch to Forgot Password ──
  const goToForgotPassword = () => {
    setView("forgot-password");
    setResetEmail("");
    setResetNewPassword("");
    setResetError("");
    setResetSuccess("");
    setError("");
  };

  // ── Switch back to Login ──
  const goToLogin = () => {
    setView("login");
    setResetError("");
    setResetSuccess("");
    setError("");
  };

  const placeholderEmail =
    role === "student"
      ? "Email (@student.unklab.ac.id)"
      : "Admin Email";

  const labelEmail =
    role === "student" ? "Email Student" : "Email Admin";

  return (
    <>
      <div className="login-wrapper">
        <div className="login-card">
          {/* ── Logo Section — UNKLAB Official Logo ──── */}
          <div className="logo-section">
            <div className="logo-img-wrapper">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src="https://upload.wikimedia.org/wikipedia/id/5/5a/Logo_Universitas_Klabat.png"
                alt="Logo UNKLAB"
                className="logo-img"
              />
            </div>
            <h1 className="logo-title">Kontrak Matakuliah</h1>
            <p className="logo-subtitle">FILKOM — Universitas Klabat</p>
          </div>

          {/* ── Credential Hint (Testing Only) ─────────── */}
          <div className="credential-hint">
            <div className="credential-hint-title">🔑 Demo Credentials</div>
            <code>
              Student: 22010001@student.unklab.ac.id / password123
              <br />
              Admin: admin@unklab.ac.id / admin123
            </code>
          </div>

          {/* ═══════════════ LOGIN VIEW ═══════════════ */}
          {view === "login" && (
            <>
              {/* ── Role Switcher ──────────────────────── */}
              <div className="role-switcher">
                <button
                  type="button"
                  id="role-student-btn"
                  className={`role-btn ${role === "student" ? "role-btn--active" : ""}`}
                  onClick={() => {
                    setRole("student");
                    setUsername("");
                    setPassword("");
                    setError("");
                  }}
                >
                  🎓 Student
                </button>
                <button
                  type="button"
                  id="role-admin-btn"
                  className={`role-btn ${role === "admin" ? "role-btn--active" : ""}`}
                  onClick={() => {
                    setRole("admin");
                    setUsername("");
                    setPassword("");
                    setError("");
                  }}
                >
                  🛡️ Admin
                </button>
              </div>

              {/* ── Login Form ─────────────────────────── */}
              <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <label htmlFor="username-input" className="form-label">
                    {labelEmail}
                  </label>
                  <input
                    id="username-input"
                    type="email"
                    className="form-input"
                    placeholder={placeholderEmail}
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                    autoComplete="email"
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="password-input" className="form-label">
                    Password
                  </label>
                  <div className="password-wrapper">
                    <input
                      id="password-input"
                      type={showPassword ? "text" : "password"}
                      className="form-input"
                      placeholder="Masukkan password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      required
                      autoComplete="current-password"
                    />
                    <button
                      type="button"
                      className="password-toggle"
                      onClick={() => setShowPassword(!showPassword)}
                      aria-label={showPassword ? "Sembunyikan password" : "Tampilkan password"}
                    >
                      {showPassword ? (
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
                          <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                          <line x1="1" y1="1" x2="23" y2="23" />
                        </svg>
                      ) : (
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                          <circle cx="12" cy="12" r="3" />
                        </svg>
                      )}
                    </button>
                  </div>
                </div>

                {/* ── Additional Links ──────────────── */}
                <div className="additional-links">
                  <button
                    type="button"
                    className="link-btn"
                    onClick={() => alert("Fitur pendaftaran (SignUp) belum tersedia di versi demo ini.")}
                  >
                    Belum punya akun? Daftar
                  </button>
                  <button
                    type="button"
                    className="link-btn"
                    onClick={goToForgotPassword}
                  >
                    Lupa Password?
                  </button>
                </div>

                <button
                  id="login-btn"
                  type="submit"
                  className="login-btn"
                  disabled={loading}
                >
                  {loading ? (
                    <>
                      <span className="spinner" />
                      Memproses...
                    </>
                  ) : (
                    "Masuk"
                  )}
                </button>

                {error && (
                  <div id="error-message" className="error-message">
                    {error}
                  </div>
                )}
              </form>
            </>
          )}

          {/* ═══════════════ FORGOT PASSWORD VIEW ═══════════════ */}
          {view === "forgot-password" && (
            <div className="forgot-password-section">
              <div className="forgot-header">
                <svg className="forgot-header-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
                <h2 className="forgot-title">Ubah Password</h2>
                <p className="forgot-subtitle">Masukkan email terdaftar dan password baru Anda</p>
              </div>

              <form onSubmit={handleResetPassword}>
                <div className="form-group">
                  <label htmlFor="reset-email-input" className="form-label">
                    Email
                  </label>
                  <input
                    id="reset-email-input"
                    type="email"
                    className="form-input"
                    placeholder="Email terdaftar (@unklab.ac.id)"
                    value={resetEmail}
                    onChange={(e) => setResetEmail(e.target.value)}
                    required
                    autoComplete="email"
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="reset-password-input" className="form-label">
                    Password Baru
                  </label>
                  <div className="password-wrapper">
                    <input
                      id="reset-password-input"
                      type={resetShowPassword ? "text" : "password"}
                      className="form-input"
                      placeholder="Minimal 6 karakter"
                      value={resetNewPassword}
                      onChange={(e) => setResetNewPassword(e.target.value)}
                      required
                      minLength={6}
                      autoComplete="new-password"
                    />
                    <button
                      type="button"
                      className="password-toggle"
                      onClick={() => setResetShowPassword(!resetShowPassword)}
                      aria-label={resetShowPassword ? "Sembunyikan password" : "Tampilkan password"}
                    >
                      {resetShowPassword ? (
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
                          <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                          <line x1="1" y1="1" x2="23" y2="23" />
                        </svg>
                      ) : (
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                          <circle cx="12" cy="12" r="3" />
                        </svg>
                      )}
                    </button>
                  </div>
                </div>

                <button
                  id="reset-password-btn"
                  type="submit"
                  className="login-btn"
                  disabled={resetLoading}
                >
                  {resetLoading ? (
                    <>
                      <span className="spinner" />
                      Memproses...
                    </>
                  ) : (
                    "Ubah Password"
                  )}
                </button>

                {resetError && (
                  <div id="reset-error-message" className="error-message">
                    {resetError}
                  </div>
                )}

                {resetSuccess && (
                  <div id="reset-success-message" className="success-message">
                    ✅ {resetSuccess}
                  </div>
                )}
              </form>

              <button
                type="button"
                id="back-to-login-btn"
                className="back-to-login-btn"
                onClick={goToLogin}
              >
                ← Kembali ke Login
              </button>
            </div>
          )}

          {/* ── Footer ─────────────────────────────────── */}
          <div className="login-footer">
            <p>
              Sistem Kontrak Matakuliah —{" "}
              <span className="campus-name">Universitas Klabat</span>
              <br />
              Fakultas Ilmu Komputer (FILKOM)
            </p>
          </div>
        </div>
      </div>

      {/* ── Success Overlay ────────────────────────── */}
      {success && (
        <div className="success-overlay">
          <div className="success-card">
            <div className="success-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                <polyline points="20 6 9 17 4 12" />
              </svg>
            </div>
            <h2 className="success-title">Login Sukses!</h2>
            <p className="success-text">Selamat datang, {loggedUser}</p>
          </div>
        </div>
      )}
    </>
  );
}
