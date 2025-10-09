"use server";

import { AgentSpec, BaseResponse } from "@/types";
import { Agent, AgentResponse, Tool } from "@/types";
import { revalidatePath } from "next/cache";
import { fetchApi, createErrorResponse } from "./utils";
import { AgentFormData } from "@/components/AgentsProvider";
import { isMcpTool } from "@/lib/toolUtils";
import { k8sRefUtils } from "@/lib/k8sUtils";

/**
 * Converts AgentFormData to Agent format
 * @param agentFormData The form data to convert
 * @returns An Agent object
 */
function fromAgentFormDataToAgent(agentFormData: AgentFormData): Agent {
  const modelConfigName = agentFormData.modelName?.includes("/")
    ? agentFormData.modelName.split("/").pop() || ""
    : agentFormData.modelName;

  const type = agentFormData.type || "Declarative";

  const convertTools = (tools: Tool[]) =>
    tools.map((tool) => {
      if (isMcpTool(tool)) {
        const mcpServer = tool.mcpServer;
        if (!mcpServer) {
          throw new Error("MCP server not found");
        }
        // Ensure TypedLocalReference fields are only the name (no namespace)
        let name = mcpServer.name;
        if (k8sRefUtils.isValidRef(mcpServer.name)) {
          name = k8sRefUtils.fromRef(mcpServer.name).name;
        }

        let kind = mcpServer.kind;
        // Special handling for kagent-querydoc - always ensure correct apiGroup
        if (mcpServer.name.toLocaleLowerCase().includes("kagent-querydoc")) {
          mcpServer.apiGroup = "";
          kind = "Service";
        }

        return {
          type: "McpServer",
          mcpServer: {
            name,
            kind,
            apiGroup: mcpServer.apiGroup,
            toolNames: mcpServer.toolNames,
          },
        } as Tool;
      }

      if (tool.type === "Agent") {
        const agentObj = tool.agent as { ref?: string; name?: string; kind?: string; apiGroup?: string };
        const refOrName = agentObj.ref || agentObj.name || "";
        const nameOnly = k8sRefUtils.isValidRef(refOrName) ? k8sRefUtils.fromRef(refOrName).name : refOrName;
        return {
          type: "Agent",
          agent: {
            name: nameOnly,
            kind: agentObj.kind || "Agent",
            apiGroup: agentObj.apiGroup || "kagent.dev",
          },
        } as Tool;
      }

      console.warn("Unknown tool type:", tool);
      return tool as Tool;
    });

  const base: Partial<Agent> = {
    metadata: {
      name: agentFormData.name,
      namespace: agentFormData.namespace || "",
    },
    spec: {
      type,
      description: agentFormData.description,
    } as AgentSpec,
  };

  if (type === "Declarative") {
    base.spec!.declarative = {
      systemMessage: agentFormData.systemPrompt || "",
      modelConfig: modelConfigName || "",
      stream: agentFormData.stream ?? true,
      tools: convertTools(agentFormData.tools || []),
      dspyConfig: agentFormData.dspyConfig,
    };
  } else if (type === "BYO") {
    base.spec!.byo = {
      deployment: {
        image: agentFormData.byoImage || "",
        cmd: agentFormData.byoCmd,
        args: agentFormData.byoArgs,
        replicas: agentFormData.replicas,
        imagePullSecrets: agentFormData.imagePullSecrets,
        volumes: agentFormData.volumes,
        volumeMounts: agentFormData.volumeMounts,
        labels: agentFormData.labels,
        annotations: agentFormData.annotations,
        env: agentFormData.env,
        imagePullPolicy: agentFormData.imagePullPolicy,
      },
    };
  }

  return base as Agent;
}

export async function getAgent(agentName: string, namespace: string): Promise<BaseResponse<AgentResponse>> {
  try {
    const agentData = await fetchApi<BaseResponse<AgentResponse>>(`/agents/${namespace}/${agentName}`);
    return { message: "Successfully fetched agent", data: agentData.data };
  } catch (error) {
    return createErrorResponse<AgentResponse>(error, "Error getting agent");
  }
}

/**
 * Deletes a agent
 * @param agentName The agent name
 * @param namespace The agent namespace
 * @returns A promise with the delete result
 */
export async function deleteAgent(agentName: string, namespace: string): Promise<BaseResponse<void>> {
  try {
    await fetchApi(`/agents/${namespace}/${agentName}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    revalidatePath("/");
    return { message: "Successfully deleted agent" };
  } catch (error) {
    return createErrorResponse<void>(error, "Error deleting agent");
  }
}

/**
 * Creates or updates an agent
 * @param agentConfig The agent configuration
 * @param update Whether to update an existing agent
 * @returns A promise with the created/updated agent
 */
export async function createAgent(agentConfig: AgentFormData, update: boolean = false): Promise<BaseResponse<Agent>> {
  try {
    // Only get the name of the model, not the full ref
    if (agentConfig.modelName) {
      if (k8sRefUtils.isValidRef(agentConfig.modelName)) {
        agentConfig.modelName = k8sRefUtils.fromRef(agentConfig.modelName).name;
      }
    }

    const agentPayload = fromAgentFormDataToAgent(agentConfig);
    const response = await fetchApi<BaseResponse<Agent>>(`/agents`, {
      method: update ? "PUT" : "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(agentPayload),
    });

    if (!response) {
      throw new Error("Failed to create agent");
    }

    const agentRef = k8sRefUtils.toRef(
      response.data!.metadata.namespace || "",
      response.data!.metadata.name,
    )

    revalidatePath(`/agents/${agentRef}/chat`);
    return { message: "Successfully created agent", data: response.data };
  } catch (error) {
    return createErrorResponse<Agent>(error, "Error creating agent");
  }
}

/**
 * Gets all agents
 * @returns A promise with all agents
 */
export async function getAgents(): Promise<BaseResponse<AgentResponse[]>> {
  try {
    const { data } = await fetchApi<BaseResponse<AgentResponse[]>>(`/agents`);

    const sortedData = data?.sort((a, b) => {
      const aRef = k8sRefUtils.toRef(a.agent.metadata.namespace || "", a.agent.metadata.name);
      const bRef = k8sRefUtils.toRef(b.agent.metadata.namespace || "", b.agent.metadata.name);
      return aRef.localeCompare(bRef);
    });

    return { message: "Successfully fetched agents", data: sortedData };
  } catch (error) {
    return createErrorResponse<AgentResponse[]>(error, "Error getting agents");
  }
}
