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
  const accessToken = useAuthStore((state) => state.accessToken);
  const queryClient = new QueryClient();

  useEffect(() => {
    const { data: listener } = supabase.auth.onAuthStateChange(
      (_event, session) => {
        if (session) {
          useAuthStore.getState().login(session.user.id, session.access_token!);
        } else {
          useAuthStore.getState().logout();
        }
      }
    );

    supabase.auth.getSession().then(({ data }) => {
      if (data.session) {
        useAuthStore
          .getState()
          .login(data.session.user.id, data.session.access_token!);
      }
    });

    return () => listener.subscription.unsubscribe();
  }, []);

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
