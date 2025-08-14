import { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { supabase } from "./lib/supabaseClient";
import { useAuthStore } from "./store/AuthStore";
import Dashboard from "./pages/Dashboard";
import DashboardV2 from "./pages/DashboardV2";
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
    const restoreSession = async () => {
      const {
        data: { session },
      } = await supabase.auth.getSession();

      if (session?.user?.id && session?.access_token) {
        login(session.user.id, session.access_token);
      }

      setSessionChecked(true);
    };

    restoreSession();

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

  if (!sessionChecked) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p>Loading session...</p>
        </div>
      </div>
    );
  }

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
