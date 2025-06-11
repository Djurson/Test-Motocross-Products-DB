"use client";

import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Label } from "./ui/label";
import { UserInput } from "@/utils/types";

type Option = {
  value: string;
  label: string;
};

const frameworks = [
  {
    value: "next.js",
    label: "Next.js",
  },
  {
    value: "sveltekit",
    label: "SvelteKit",
  },
  {
    value: "nuxt.js",
    label: "Nuxt.js",
  },
  {
    value: "remix",
    label: "Remix",
  },
  {
    value: "astro",
    label: "Astro",
  },
];

export function DropDown({
  label,
  placeholder,
  disabled,
  input,
  setInput,
  type,
}: {
  label: string;
  placeholder: string;
  disabled: boolean;
  input: UserInput;
  setInput: React.Dispatch<React.SetStateAction<UserInput>>;
  type: keyof UserInput;
}) {
  const [open, setOpen] = React.useState(false);

  const selectedValue = input[type];
  const selectedString =
    typeof selectedValue === "string"
      ? selectedValue
      : typeof selectedValue === "object" && selectedValue !== null
      ? "name" in selectedValue
        ? selectedValue.name
        : "sizeCC" in selectedValue
        ? String(selectedValue.sizeCC)
        : "category" in selectedValue
        ? selectedValue.category
        : undefined
      : undefined;

  return (
    <div className="flex flex-col items-start justify-center gap-1">
      <Label className="text-sm">{label}</Label>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild disabled={disabled}>
          <Button variant="outline" role="combobox" aria-expanded={open} className="w-[200px] justify-between">
            {selectedString ? frameworks.find((f) => f.value === selectedString)?.label : placeholder}
            <ChevronsUpDown className="opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0">
          <Command>
            <CommandInput placeholder={placeholder} className="h-9" disabled={disabled} />
            <CommandList>
              <CommandEmpty>Inget {label} hittat</CommandEmpty>
              <CommandGroup>
                {frameworks.map((framework) => (
                  <CommandItem
                    key={framework.value}
                    value={framework.value}
                    onSelect={(currentValue) => {
                      setInput({
                        ...input,
                        [type]: currentValue,
                      });
                      setOpen(false);
                    }}>
                    {framework.label}
                    <Check className={cn("ml-auto", selectedString === framework.value ? "opacity-100" : "opacity-0")} />
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  );
}
