import { useNavigate } from "react-router-dom";
import { supabase } from "../lib/supabaseClient";
import { useAuthStore } from "../store/AuthStore";

export const SignOutButton = () => {
  const logout = useAuthStore((state) => state.logout);
  const navigate = useNavigate();

  const handleSignOut = async () => {
    await supabase.auth.signOut();
    logout();
    navigate("/");
  };

  return (
    <button
      className="mt-2 px-4 py-2 rounded-md font-medium text-indigo-500 hover:text-white transition-all duration-200 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 shadow-sm hover:shadow-md bg-red-100 hover:bg-indigo-500 border border-indigo-500"
      onClick={handleSignOut}
    >
      Sign Out
    </button>
  );
};
