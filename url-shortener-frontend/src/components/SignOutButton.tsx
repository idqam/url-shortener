// src/components/SignOutButton.tsx
import { useNavigate } from "react-router-dom";
import { supabase } from "../lib/supabaseClient"; // adjust if needed
import { useAuthStore } from "../store/AuthStore";

export const SignOutButton = () => {
  const logout = useAuthStore((state) => state.logout);
  const navigate = useNavigate();

  const handleSignOut = async () => {
    await supabase.auth.signOut(); // ends Supabase session
    logout(); // clear Zustand store
    navigate("/"); // redirect to homepage
  };

  return <button onClick={handleSignOut}>Sign Out</button>;
};
