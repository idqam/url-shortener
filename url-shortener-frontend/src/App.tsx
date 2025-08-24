import { useEffect } from "react";
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

  useEffect(() => {
    const { data: listener } = supabase.auth.onAuthStateChange(
      (_event, session) => {
        if (session?.user?.id && session?.access_token) {
          login(session.user.id, session.access_token);
        } else {
          logout();
        }
      }
    );

    supabase.auth.getSession().then(({ data, error }) => {
      if (error) {
        console.error("Session error:", error);
        logout();
      } else if (data.session?.user?.id && data.session?.access_token) {
        login(data.session.user.id, data.session.access_token);
      }
    });

    return () => listener.subscription.unsubscribe();
  }, [login, logout]);

  return (
    <QueryClientProvider client={queryClient}>
      <Router>
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
