"use client";

import { RefreshCw, Clock } from "lucide-react";
import { useAgents } from "./AgentsProvider";
import { cn } from "@/lib/utils";

export function RefreshIndicator() {
  const { isRefreshing, lastRefresh } = useAgents();

  if (!isRefreshing && !lastRefresh) {
    return null;
  }

  const formatLastRefresh = (date: Date) => {
    const now = new Date();
    const diff = Math.floor((now.getTime() - date.getTime()) / 1000);
    
    if (diff < 60) {
      return "just now";
    } else if (diff < 3600) {
      const minutes = Math.floor(diff / 60);
      return `${minutes}m ago`;
    } else {
      const hours = Math.floor(diff / 3600);
      return `${hours}h ago`;
    }
  };

  return (
    <div className="flex items-center gap-2 text-sm text-muted-foreground">
      {isRefreshing ? (
        <>
          <RefreshCw className={cn("h-4 w-4 animate-spin")} />
          <span>Updating...</span>
        </>
      ) : lastRefresh ? (
        <>
          <Clock className="h-4 w-4" />
          <span>Updated {formatLastRefresh(lastRefresh)}</span>
        </>
      ) : null}
    </div>
  );
}