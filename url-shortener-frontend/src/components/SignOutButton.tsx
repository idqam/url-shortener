import { useAuthStore } from "../store/AuthStore";
import { supabase } from "../lib/supabaseClient";

export default function SignOutButton() {
  const logout = useAuthStore((state) => state.logout);

  const handleSignOut = async () => {
    await supabase.auth.signOut();
    logout();
  };

  return (
    <button
      onClick={handleSignOut}
      className="text-sm font-medium text-gray-600 hover:text-gray-900 px-3 py-1.5 rounded-md hover:bg-gray-100 transition-colors"
    >
      Sign out
    </button>
  );
}
