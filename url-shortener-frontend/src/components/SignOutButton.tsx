import { useAuthStore } from "../store/AuthStore";
import { supabase } from "../lib/supabaseClient";
import { analyticsColors } from "../constants/colors";

export default function SignOutButton() {
  const logout = useAuthStore((state) => state.logout);

  const handleSignOut = async () => {
    await supabase.auth.signOut();
    logout();
  };

  return (
    <button
      onClick={handleSignOut}
      className="mt-6 px-5 py-2 rounded-xl text-sm font-semibold text-white transition-all duration-300 hover:scale-105"
      style={{
        backgroundImage: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
      }}
    >
      Sign Out
    </button>
  );
}
