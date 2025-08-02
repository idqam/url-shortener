/* eslint-disable @typescript-eslint/no-explicit-any */
import { useUrlStore } from "../store/UrlStore";
import { shortenAuthed } from "../api/urls";
import type { UrlRequestAuthed, UrlStructure } from "../types/urls";

export function useShortenUrlAuth() {
  const setLoading = useUrlStore((s) => s.setLoading);
  const setError = useUrlStore((s) => s.setError);
  const addAuthedUrl = useUrlStore((s) => s.addAuthedUrl);

  const shorten = async (payload: UrlRequestAuthed) => {
    setLoading(true);
    setError(null);
    try {
      const result = await shortenAuthed(payload);

      const urlStructure: UrlStructure = {
        id: result.id,
        original_url: result.original_url,
        short_url: result.short_url,
        short_code: result.short_code,
        user_id: payload.user_id,
        click_count: 0,
        created_at: new Date().toISOString(),
        is_public: payload.is_public,
      };

      addAuthedUrl(urlStructure);
    } catch (err: any) {
      setError(err.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  return { shorten };
}
