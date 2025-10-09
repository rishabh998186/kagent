"use client";
import { DSPyConfig } from "@/types";
import React from "react";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { SystemPromptSection } from "./SystemPromptSection";
import { DSPyConfigSection } from "./DSPyConfigSection";

interface PromptConfigSectionProps {
  // Traditional system prompt props
  systemPrompt: string;
  onSystemPromptChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  onSystemPromptBlur?: () => void;
  systemPromptError?: string;
  
  // DSPy config props
  useDSPy: boolean;
  onDSPyToggle: (enabled: boolean) => void;
  dspyConfig: DSPyConfig;
  onDSPyConfigChange: (config: DSPyConfig) => void;
  
  disabled?: boolean;
}

export function PromptConfigSection({
  systemPrompt,
  onSystemPromptChange,
  onSystemPromptBlur,
  systemPromptError,
  useDSPy,
  onDSPyToggle,
  dspyConfig,
  onDSPyConfigChange,
  disabled,
}: PromptConfigSectionProps) {
  return (
    <div className="space-y-6">
      {/* Toggle between traditional prompt and DSPy */}
      <div className="flex items-center justify-between p-4 border rounded-lg bg-muted/50">
        <div>
          <Label className="text-base font-bold">Prompt Mode</Label>
          <p className="text-xs text-muted-foreground mt-1">
            Choose between traditional system prompt or DSPy-based prompt compilation
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Label htmlFor="dspy-toggle" className="text-sm font-medium">
            {useDSPy ? "DSPy Mode" : "Traditional Mode"}
          </Label>
          <Switch
            id="dspy-toggle"
            checked={useDSPy}
            onCheckedChange={onDSPyToggle}
            disabled={disabled}
          />
        </div>
      </div>

      {/* Render the appropriate section */}
      {!useDSPy ? (
        <SystemPromptSection
          value={systemPrompt}
          onChange={onSystemPromptChange}
          onBlur={onSystemPromptBlur}
          error={systemPromptError}
          disabled={disabled || false}
        />
      ) : (
        <DSPyConfigSection
          config={dspyConfig}
          onChange={onDSPyConfigChange}
          disabled={disabled}
        />
      )}
    </div>
  );
}
