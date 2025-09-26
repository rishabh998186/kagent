"use server";

import { BaseResponse, ToolsResponse } from "@/types";
import { fetchApi } from "./utils";

/**
 * Gets all available tools
 * @returns A promise with all tools
 */
export async function getTools(): Promise<ToolsResponse[]> {
  try {
    const response = await fetchApi<BaseResponse<ToolsResponse[]>>("/tools");
    if (!response) {
      console.warn("No response received from tools API");
      return [];
    }
    if (response.error) {
      console.warn("Tools API returned error:", response.error);
      return [];
    }
    return response.data || [];
  } catch (error) {
    console.error("Error getting built-in tools:", error);
    // Return empty array instead of throwing to prevent UI crashes
    return [];
  }
}
