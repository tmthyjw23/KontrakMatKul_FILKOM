"use client";

import {
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { Toaster } from "sonner";
import { useAuthStore } from "@/lib/store/useAuthStore";

type ProvidersProps = {
  children: React.ReactNode;
};

function AuthBootstrap({ children }: { children: React.ReactNode }) {
  const loadFromStorage = useAuthStore((state) => state.loadFromStorage);

  useEffect(() => {
    loadFromStorage();
  }, [loadFromStorage]);

  return <>{children}</>;
}

export function Providers({ children }: ProvidersProps) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 5 * 60 * 1000,
            gcTime: 30 * 60 * 1000,
            retry: 1,
            refetchOnWindowFocus: false,
          },
          mutations: {
            retry: 0,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      <AuthBootstrap>
        {children}
        <Toaster
          position="top-right"
          richColors={false}
          toastOptions={{
            className: "!border !border-white/10 !bg-zinc-950/90 !text-zinc-100",
          }}
        />
      </AuthBootstrap>
    </QueryClientProvider>
  );
}
