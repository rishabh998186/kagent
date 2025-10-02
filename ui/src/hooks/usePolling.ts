import { useRef, useEffect, useCallback } from 'react';

interface UsePollingOptions {
  interval?: number; // Polling interval in milliseconds
  immediate?: boolean; // Whether to call the function immediately
  enabled?: boolean; // Whether polling is enabled
  onlyWhenVisible?: boolean; // Only poll when page is visible
}

/**
 * Custom hook for polling data at regular intervals
 * @param callback The function to call repeatedly
 * @param options Polling configuration options
 */
export function usePolling(
  callback: () => void | Promise<void>,
  options: UsePollingOptions = {}
) {
  const {
    interval = 30000, // Default 30 seconds
    immediate = false,
    enabled = true,
    onlyWhenVisible = true
  } = options;

  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const callbackRef = useRef(callback);
  const isPollingRef = useRef(false);

  // Keep callback reference updated
  useEffect(() => {
    callbackRef.current = callback;
  }, [callback]);

  const startPolling = useCallback(() => {
    if (isPollingRef.current || !enabled) return;
    
    isPollingRef.current = true;
    
    const poll = async () => {
      // Check if page is visible (if onlyWhenVisible is enabled)
      if (onlyWhenVisible && document.hidden) {
        return;
      }
      
      try {
        await callbackRef.current();
      } catch (error) {
        console.error('Polling error:', error);
      }
    };

    // Call immediately if requested
    if (immediate) {
      poll();
    }

    // Set up interval
    intervalRef.current = setInterval(poll, interval);
  }, [interval, immediate, enabled, onlyWhenVisible]);

  const stopPolling = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
    isPollingRef.current = false;
  }, []);

  // Handle visibility change
  useEffect(() => {
    if (!onlyWhenVisible) return;

    const handleVisibilityChange = () => {
      if (document.hidden) {
        stopPolling();
      } else if (enabled) {
        startPolling();
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [onlyWhenVisible, enabled, startPolling, stopPolling]);

  // Start/stop polling based on enabled state
  useEffect(() => {
    if (enabled) {
      startPolling();
    } else {
      stopPolling();
    }

    return stopPolling;
  }, [enabled, startPolling, stopPolling]);

  return {
    startPolling,
    stopPolling,
    isPolling: isPollingRef.current
  };
}