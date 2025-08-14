import { useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Home } from "./pages/Home";
import { SignUpPage } from "./pages/SignUpPage";
import LogInPage from "./pages/LogInPage";
import { supabase } from "./lib/supabaseClient";
import { useAuthStore } from "./store/AuthStore";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import Dashboard from "./pages/Dashboard";
import { ProtectedRoute } from "./components/ProtectedRoute";
import DashboardV2 from "./pages/DashboardV2";

function App() {
  const login = useAuthStore((state) => state.login);
  const logout = useAuthStore((state) => state.logout);
  const accessToken = useAuthStore((state) => state.accessToken);
  const queryClient = new QueryClient();

  useEffect(() => {
    const { data: listener } = supabase.auth.onAuthStateChange(
      async (event, session) => {
        if (event === "SIGNED_IN" && session?.user && session?.access_token) {
          login(session.user.id, session.access_token);
        }
        if (event === "SIGNED_OUT") {
          logout();
        }
      }
    );

    return () => {
      listener.subscription.unsubscribe();
    };
  }, [login, logout]);

  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/signup" element={<SignUpPage />} />
          <Route path="/login" element={<LogInPage />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard token={accessToken as string} />
              </ProtectedRoute>
            }
          />
          <Route
            path="/dashboard2"
            element={
              <ProtectedRoute>
                <DashboardV2 token={accessToken as string} />
              </ProtectedRoute>
            }
          />
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}

export default App;
