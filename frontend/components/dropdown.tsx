"use client";

import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Label } from "./ui/label";
import { UserInput } from "@/utils/types";

type DropDownProps<T> = {
  label: string;
  placeholder: string;
  disabled?: boolean;
  input: UserInput;
  setInput: React.Dispatch<React.SetStateAction<UserInput>>;
  type: keyof UserInput;
  data: T[];
  getOptionLabel: (item: T) => string;
  getOptionValue: (item: T) => string;
};

export function DropDown<T>({ label, placeholder, disabled, input, setInput, type, data, getOptionLabel, getOptionValue }: DropDownProps<T>) {
  const [open, setOpen] = React.useState(false);
  const selectedItem = input[type] as T | undefined;
  const selectedValue = selectedItem ? getOptionValue(selectedItem) : undefined;

  return (
    <div className="flex flex-col items-start justify-center w-full gap-1">
      <Label className="text-sm">{label}</Label>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild disabled={disabled}>
          <Button variant="outline" role="combobox" aria-expanded={open} className="w-[200px] justify-between">
            {selectedValue ? getOptionLabel(data.find((d) => getOptionValue(d) === selectedValue) as T) : placeholder}
            <ChevronsUpDown className="w-4 h-4 ml-2 opacity-50" />
          </Button>
        </PopoverTrigger>
        {data && (
          <PopoverContent className="w-[200px] p-0">
            <Command>
              <CommandInput placeholder={placeholder} className="h-9" disabled={disabled} />
              <CommandList>
                <CommandEmpty>Inget {label} hittat</CommandEmpty>
                <CommandGroup>
                  {data.map((item) => {
                    const value = getOptionValue(item);
                    const label = getOptionLabel(item);
                    return (
                      <CommandItem
                        key={value}
                        value={value}
                        onSelect={() => {
                          let updatedInput: UserInput = { ...input, [type]: item };

                          if (type === "brand") {
                            updatedInput.model = undefined;
                            updatedInput.year = undefined;
                          } else if (type === "model") {
                            updatedInput.year = undefined;
                          }

                          setInput(updatedInput);
                          setOpen(false);
                        }}>
                        {label}
                        <Check className={cn("ml-auto", selectedValue === value ? "opacity-100" : "opacity-0")} />
                      </CommandItem>
                    );
                  })}
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        )}
      </Popover>
    </div>
  );
}
