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

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5,
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  const accessToken = useAuthStore((state) => state.accessToken);

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

    const initializeSession = async () => {
      try {
        const { data } = await supabase.auth.getSession();
        if (data.session) {
          useAuthStore
            .getState()
            .login(data.session.user.id, data.session.access_token!);
        }
      } catch (error) {
        console.error("Error initializing session:", error);
      }
    };

    initializeSession();

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
