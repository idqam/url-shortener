import { Link as LinkIcon } from "lucide-react";
import { useAuthStore } from "../store/AuthStore";
import { supabase } from "../lib/supabaseClient";

interface NavbarProps {
  onToggleShortener: () => void;
  showShortener: boolean;
}

export function Navbar({ onToggleShortener, showShortener }: NavbarProps) {
  const logout = useAuthStore((state) => state.logout);

  const handleSignOut = async () => {
    await supabase.auth.signOut();
    logout();
  };

  return (
    <header className="sticky top-0 z-10 bg-white border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-6 h-14 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-7 h-7 bg-blue-600 rounded-md flex items-center justify-center">
            <LinkIcon className="w-4 h-4 text-white" />
          </div>
          <span className="font-semibold text-gray-900 text-sm">Shortlink</span>
        </div>
        <div className="flex items-center gap-1">
          <button
            onClick={onToggleShortener}
            className="text-sm font-medium text-gray-600 hover:text-gray-900 px-3 py-1.5 rounded-md hover:bg-gray-100 transition-colors"
          >
            {showShortener ? "Hide Shortener" : "Quick Shorten"}
          </button>
          <button
            onClick={handleSignOut}
            className="text-sm font-medium text-gray-600 hover:text-gray-900 px-3 py-1.5 rounded-md hover:bg-gray-100 transition-colors"
          >
            Sign out
          </button>
        </div>
      </div>
    </header>
  );
}
