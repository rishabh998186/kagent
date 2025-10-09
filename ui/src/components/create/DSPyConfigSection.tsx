"use client";
import React from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { PlusCircle, Trash2, Sparkles } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { DSPyConfig, DSPyField } from "@/types";

interface DSPyConfigSectionProps {
  config: DSPyConfig;
  onChange: (config: DSPyConfig) => void;
  disabled?: boolean;
}

export function DSPyConfigSection({ config, onChange, disabled }: DSPyConfigSectionProps) {
  const updateConfig = (updates: Partial<DSPyConfig>) => {
    onChange({ ...config, ...updates });
  };

  const updateSignature = (updates: Partial<DSPyConfig["signature"]>) => {
    updateConfig({ signature: { ...config.signature, ...updates } });
  };

  const addField = (type: "inputs" | "outputs") => {
    const newField: DSPyField = { name: "", type: "string", description: "" };
    updateSignature({
      [type]: [...config.signature[type], newField],
    });
  };

  const removeField = (type: "inputs" | "outputs", index: number) => {
    updateSignature({
      [type]: config.signature[type].filter((_, i) => i !== index),
    });
  };

  const updateField = (type: "inputs" | "outputs", index: number, field: DSPyField) => {
    const updated = [...config.signature[type]];
    updated[index] = field;
    updateSignature({ [type]: updated });
  };

  const renderFieldList = (type: "inputs" | "outputs") => {
    const fields = config.signature[type];
    const label = type === "inputs" ? "Input Fields" : "Output Fields";

    return (
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <Label className="text-sm font-semibold">{label}</Label>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => addField(type)}
            disabled={disabled}
          >
            <PlusCircle className="h-4 w-4 mr-1" />
            Add {type === "inputs" ? "Input" : "Output"}
          </Button>
        </div>

        {fields.map((field, index) => (
          <Card key={index} className="p-3">
            <div className="space-y-3">
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <Label className="text-xs">Field Name</Label>
                  <Input
                    placeholder="e.g., question, context"
                    value={field.name}
                    onChange={(e) =>
                      updateField(type, index, { ...field, name: e.target.value })
                    }
                    disabled={disabled}
                    className="h-9"
                  />
                </div>
                <div>
                  <Label className="text-xs">Field Type</Label>
                  <Select
                    value={field.type}
                    onValueChange={(val) =>
                      updateField(type, index, { ...field, type: val as DSPyField["type"] })
                    }
                    disabled={disabled}
                  >
                    <SelectTrigger className="h-9">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="string">String</SelectItem>
                      <SelectItem value="int">Integer</SelectItem>
                      <SelectItem value="float">Float</SelectItem>
                      <SelectItem value="bool">Boolean</SelectItem>
                      <SelectItem value="list">List</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div>
                <Label className="text-xs">Description</Label>
                <Input
                  placeholder="Describe this field..."
                  value={field.description || ""}
                  onChange={(e) =>
                    updateField(type, index, { ...field, description: e.target.value })
                  }
                  disabled={disabled}
                  className="h-9"
                />
              </div>

              <div className="flex justify-end">
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => removeField(type, index)}
                  disabled={disabled || fields.length === 1}
                  className="text-red-500 hover:text-red-700"
                >
                  <Trash2 className="h-4 w-4 mr-1" />
                  Remove
                </Button>
              </div>
            </div>
          </Card>
        ))}

        {fields.length === 0 && (
          <p className="text-sm text-muted-foreground text-center py-4">
            No {type} defined. Click &quot;Add {type === "inputs" ? "Input" : "Output"}&quot; to get started.
          </p>
        )}
      </div>
    );
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Sparkles className="h-5 w-5 text-purple-500" />
            DSPy Configuration
          </CardTitle>
          <CardDescription>
            Configure DSPy-based prompt compilation and optimization for your agent
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Module Selection */}
          <div>
            <Label className="text-sm font-semibold mb-2 block">DSPy Module</Label>
            <p className="text-xs text-muted-foreground mb-2">
              Choose the DSPy reasoning module to use
            </p>
            <Select
              value={config.module}
              onValueChange={(val) => updateConfig({ module: val as DSPyConfig["module"] })}
              disabled={disabled}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="Predict">Predict - Direct prediction</SelectItem>
                <SelectItem value="ChainOfThought">Chain of Thought - Step-by-step reasoning</SelectItem>
                <SelectItem value="ReAct">ReAct - Reasoning + Acting</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Instructions */}
          <div>
            <Label className="text-sm font-semibold mb-2 block">Instructions</Label>
            <p className="text-xs text-muted-foreground mb-2">
              Provide high-level instructions for what the agent should do
            </p>
            <Textarea
              placeholder="e.g., Answer questions based on the provided context..."
              value={config.signature.instructions}
              onChange={(e) => updateSignature({ instructions: e.target.value })}
              disabled={disabled}
              className="min-h-[100px]"
            />
          </div>

          {/* Input Fields */}
          {renderFieldList("inputs")}

          {/* Output Fields */}
          {renderFieldList("outputs")}

          {/* Optimization Config */}
          <Card className="border-dashed">
            <CardHeader>
              <CardTitle className="text-sm">Optimization (Optional)</CardTitle>
              <CardDescription className="text-xs">
                Enable prompt optimization using DSPy optimizers
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label className="text-xs">Optimizer Algorithm</Label>
                <Select
                  value={config.optimizationConfig?.optimizer || "MIPROv2"}
                  onValueChange={(val) =>
                    updateConfig({
                      optimizationConfig: {
                        ...(config.optimizationConfig || { enabled: true }),
                        optimizer: val as "MIPRO" | "MIPROv2" | "BootstrapFewShot" | "BootstrapFewShotWithRandomSearch" | "COPRO",
                      },
                    })
                  }
                  disabled={disabled}
                >
                  <SelectTrigger className="h-9">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="MIPROv2">MIPROv2 - Multi-stage optimization</SelectItem>
                    <SelectItem value="MIPRO">MIPRO - Basic optimization</SelectItem>
                    <SelectItem value="BootstrapFewShot">Bootstrap Few-Shot</SelectItem>
                    <SelectItem value="COPRO">COPRO - Coordinate ascent</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label className="text-xs">Metric Name (Optional)</Label>
                <Input
                  placeholder="e.g., accuracy, f1_score"
                  value={config.optimizationConfig?.metricName || ""}
                  onChange={(e) =>
                    updateConfig({
                      optimizationConfig: {
                        ...(config.optimizationConfig || { enabled: true, optimizer: "MIPROv2" }),
                        metricName: e.target.value,
                      },
                    })
                  }
                  disabled={disabled}
                  className="h-9"
                />
              </div>
            </CardContent>
          </Card>
        </CardContent>
      </Card>
    </div>
  );
}
