import { useEffect } from "react";
import { supabase } from "./supabaseClient";
//TESTING PURPOSES
export default function DebugToken() {
  useEffect(() => {
    (async () => {
      const { data, error } = await supabase.auth.getSession();
      if (error) {
        console.error("Session error:", error);
        return;
      }

      const token = data.session?.access_token;
      console.log("Supabase access token:", token);
    })();
  }, []);

  return <div>Check console for </div>;
}
