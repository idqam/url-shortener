import { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { supabase } from "./lib/supabaseClient";
import { useAuthStore } from "./store/AuthStore";
import Dashboard from "./pages/Dashboard";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { Home } from "./pages/Home";
import LogInPage from "./pages/LogInPage";
import { SignUpPage } from "./pages/SignUpPage";

function App() {
  const login = useAuthStore((state) => state.login);
  const logout = useAuthStore((state) => state.logout);
  const accessToken = useAuthStore((state) => state.accessToken);
  const queryClient = new QueryClient();

  const [sessionChecked, setSessionChecked] = useState(false);

  useEffect(() => {
    supabase.auth.getSession().then(({ data, error }) => {
      if (error) {
        console.error("Session error:", error);
        logout();
      } else {
        const session = data.session;
        if (session?.user?.id && session?.access_token) {
          login(session.user.id, session.access_token);
        } else {
          logout();
        }
      }
      setSessionChecked(true);
    });

    const { data: listener } = supabase.auth.onAuthStateChange(
      (_event, session) => {
        if (session?.user?.id && session?.access_token) {
          login(session.user.id, session.access_token);
        } else {
          logout();
        }
      }
    );

    return () => listener.subscription.unsubscribe();
  }, [login, logout]);

  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        {!sessionChecked && (
          <div className="fixed inset-0 flex items-center justify-center bg-white/70 z-50">
            <div className="text-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p>Checking session...</p>
            </div>
          </div>
        )}

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<LogInPage />} />
          <Route path="/signup" element={<SignUpPage />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard token={accessToken as string} />
              </ProtectedRoute>
            }
          />
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}

export default App;
