import { createClient } from "@supabase/supabase-js";

export const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL!,
  import.meta.env.VITE_SUPABASE_ANON_KEY!,
  {
    auth: {
      persistSession: true,
      autoRefreshToken: true,
      detectSessionInUrl: true,
    },
    global: {
      headers: {
        "X-Client-Info": "url-shortener-app",
      },
    },
  }
);

const SESSION_TIMEOUT = 2 * 60 * 60 * 1000;

export const setupSessionTimeout = () => {
  let timeoutId: number;

  const resetTimeout = () => {
    clearTimeout(timeoutId);
    timeoutId = window.setTimeout(() => {
      supabase.auth.signOut();
      alert("Session expired. Please sign in again.");
    }, SESSION_TIMEOUT);
  };

  ["mousedown", "mousemove", "keypress", "scroll", "touchstart"].forEach(
    (event) => {
      document.addEventListener(event, resetTimeout, { passive: true });
    }
  );

  resetTimeout();
};
