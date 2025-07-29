import { create } from "zustand";
import type { UrlStructure, ShortenResponse } from "../types/urls";

type UrlStoreState = {
  anonLastResult: ShortenResponse | null;
  urls: UrlStructure[];
  loading: boolean;
  error: string | null;

  setAnonResult: (url: ShortenResponse) => void;
  addAuthedUrl: (url: UrlStructure) => void;
  setLoading: (loading: boolean) => void;
  setError: (err: string | null) => void;
  clearOnLogout: () => void;
};

export const useUrlStore = create<UrlStoreState>((set) => ({
  anonLastResult: null,
  urls: [],
  loading: false,
  error: null,

  setAnonResult: (url) => set({ anonLastResult: url }),
  addAuthedUrl: (url) => set((state) => ({ urls: [url, ...state.urls] })),
  setLoading: (loading) => set({ loading }),
  setError: (err) => set({ error: err }),

  clearOnLogout: () =>
    set({
      urls: [],
      anonLastResult: null,
      loading: false,
      error: null,
    }),
}));
