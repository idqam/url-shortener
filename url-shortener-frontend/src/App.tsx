import { useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Home } from "./pages/Home";
import { SignUpPage } from "./pages/SignUpPage";
import Dashboard from "./pages/Dashboard";
import LogInPage from "./pages/LogInPage";
import { supabase } from "./lib/supabaseClient";
import { useAuthStore } from "./store/AuthStore";

function App() {
  const login = useAuthStore((state) => state.login);
  const logout = useAuthStore((state) => state.logout);

  useEffect(() => {
    const { data: listener } = supabase.auth.onAuthStateChange(
      async (event, session) => {
        if (event === "SIGNED_IN" && session?.user) {
          login(session.user.id);

          // const token = session.access_token;
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
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signup" element={<SignUpPage />} />
        <Route path="/login" element={<LogInPage />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </Router>
  );
}

export default App;
