"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";

interface SettingsContextType {
  autoRefreshEnabled: boolean;
  autoRefreshInterval: number;
  setAutoRefreshEnabled: (enabled: boolean) => void;
  setAutoRefreshInterval: (interval: number) => void;
}

const SettingsContext = createContext<SettingsContextType | undefined>(undefined);

export function useSettings() {
  const context = useContext(SettingsContext);
  if (context === undefined) {
    throw new Error("useSettings must be used within a SettingsProvider");
  }
  return context;
}

interface SettingsProviderProps {
  children: ReactNode;
}

// Default settings
const DEFAULT_AUTO_REFRESH_ENABLED = true;
const DEFAULT_AUTO_REFRESH_INTERVAL = 30000; // 30 seconds

export function SettingsProvider({ children }: SettingsProviderProps) {
  const [autoRefreshEnabled, setAutoRefreshEnabledState] = useState(DEFAULT_AUTO_REFRESH_ENABLED);
  const [autoRefreshInterval, setAutoRefreshIntervalState] = useState(DEFAULT_AUTO_REFRESH_INTERVAL);

  // Load settings from localStorage on mount
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const savedAutoRefreshEnabled = localStorage.getItem('kagent.autoRefreshEnabled');
      const savedAutoRefreshInterval = localStorage.getItem('kagent.autoRefreshInterval');

      if (savedAutoRefreshEnabled !== null) {
        setAutoRefreshEnabledState(JSON.parse(savedAutoRefreshEnabled));
      }

      if (savedAutoRefreshInterval !== null) {
        setAutoRefreshIntervalState(parseInt(savedAutoRefreshInterval, 10));
      }
    }
  }, []);

  const setAutoRefreshEnabled = (enabled: boolean) => {
    setAutoRefreshEnabledState(enabled);
    if (typeof window !== 'undefined') {
      localStorage.setItem('kagent.autoRefreshEnabled', JSON.stringify(enabled));
    }
  };

  const setAutoRefreshInterval = (interval: number) => {
    setAutoRefreshIntervalState(interval);
    if (typeof window !== 'undefined') {
      localStorage.setItem('kagent.autoRefreshInterval', interval.toString());
    }
  };

  const value = {
    autoRefreshEnabled,
    autoRefreshInterval,
    setAutoRefreshEnabled,
    setAutoRefreshInterval,
  };

  return (
    <SettingsContext.Provider value={value}>
      {children}
    </SettingsContext.Provider>
  );
}