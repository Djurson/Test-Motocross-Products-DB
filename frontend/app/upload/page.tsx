"use client";

import { ChangeEvent, DragEvent, useState } from "react";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { FileText, FileUp, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { tryCatch } from "@/utils/trycatch";
import axios from "axios";

export default function Page() {
  const [isDragging, setIsDragging] = useState(false);
  const [category, setCategory] = useState("");
  const [fileToUpload, setFileToUpload] = useState<File | undefined>();

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);

    const file = e.dataTransfer.files?.[0];
    if (!file || !file.name.toLowerCase().endsWith(".csv")) return;

    setFileToUpload(file);
  };

  async function UploadFile(e: ChangeEvent<HTMLInputElement>) {
    e.preventDefault();

    const file = e.target.files?.[0];
    if (!file || !file.name.toLowerCase().endsWith(".csv")) return;

    setFileToUpload(file);
  }

  async function ButtonClickUpload() {
    if (!fileToUpload) return;
    await UploadFileToBackend(fileToUpload, category);
  }

  return (
    <div className="flex flex-col items-center justify-center w-full gap-8 py-6">
      <Card className="w-full max-w-6xl">
        <CardHeader>
          <CardTitle>Ladda upp produkter</CardTitle>
          <CardDescription>Ladda upp produkter och koppla dem till kategorier</CardDescription>
        </CardHeader>
        <CardContent>
          <>
            {!fileToUpload ? (
              <div
                onDragOver={(e) => {
                  e.preventDefault();
                  setIsDragging(true);
                }}
                onDragLeave={(e) => {
                  e.preventDefault();
                  setIsDragging(false);
                }}
                onDrop={handleDrop}>
                <Label
                  htmlFor="PDF-Upload"
                  className={`flex flex-col items-center justify-center p-4 text-sm transition duration-300 ease-in-out border-1 
                              border-foreground border-dashed cursor-pointer rounded-2xl group shadow-[2px_4px_12px_0px_rgba(0,_0,_0,_0.08)] w-full h-full
                              ${isDragging ? "bg-accent" : "hover:bg-background bg-accent"}`}>
                  <span
                    className={`transition duration-300 ease-in-out aspect-square p-4 rounded-2xl
                                ${isDragging ? "bg-accent" : "bg-background group-hover:bg-accent"}`}>
                    <FileUp className="text-foreground aspect-square h-8.5 w-8.5" />
                  </span>
                  <span>
                    <span className="text-blue-500 transition duration-300 ease-in-out">Klicka</span> eller dra & släpp för att ladda upp produktlista
                  </span>
                  <span className="text-xs font-light text-muted-foreground">Format som stöds: .CSV</span>
                </Label>
                <Input type="file" accept=".csv" name="PDF-Upload" id="PDF-Upload" className="hidden appearance-none" onChange={UploadFile} />
              </div>
            ) : (
              <div className="w-full h-full flex flex-col items-center justify-start gap-4">
                <div className="flex items-center justify-center p-4 text-sm transition duration-300 ease-in-out border-1 gap-2 border-foreground border-dashed rounded-2xl group shadow-[2px_4px_12px_0px_rgba(0,_0,_0,_0.08)] w-full h-full bg-accent">
                  <div className="bg-background flex gap-2 w-2/3 items-center px-4 rounded-md shadow-background shadow-md">
                    <span className="transition duration-300 ease-in-out aspect-square p-4 rounded-2xl">
                      <FileText className="text-foreground aspect-square h-8.5 w-8.5" />
                    </span>
                    <div className="flex flex-col gap-2 justify-between w-full">
                      <div className="flex justify-between items-center">
                        <p className="w-fit">{fileToUpload.name}</p>
                        <p className="text-xs text-accent-foreground bg-accent px-2 py-0.5 rounded-xs">{(fileToUpload.size / (1024 * 1024)).toFixed(2)} MB</p>
                      </div>
                      <div className="flex gap-1 justify-between items-center">
                        <p className="text-xs text-accent-foreground">{fileToUpload.type}</p>
                        <p className="text-xs text-accent-foreground bg-accent px-2 py-0.5 rounded-xs">{new Date(fileToUpload.lastModified).toLocaleDateString()}</p>
                      </div>
                    </div>
                  </div>
                  <X className="relative w-6 h-6 text-red-600 cursor-pointer -top-6 -right-40 aspect-square" onClick={() => setFileToUpload(undefined)} />
                </div>
                <div className="flex justify-center items-center gap-2 w-full">
                  <div className="grid grid-cols-1 grid-rows-2 w-xs">
                    <Label>Huvud kategori</Label>
                    <Input placeholder="Exempelvis fjädring, kolvar etc..." onChange={(e) => setCategory(e.target.value)} value={category} />
                  </div>
                </div>
                <Button size="lg" onClick={ButtonClickUpload}>
                  Ladda upp
                </Button>
              </div>
            )}
          </>
        </CardContent>
      </Card>
    </div>
  );
}

async function UploadFileToBackend(file: File, rootCategory: string) {
  const formData = new FormData();

  formData.append("file", file);
  formData.append("category", rootCategory);

  const { data, error } = await tryCatch(
    axios.post("http://localhost:8000/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    })
  );

  if (error !== null) {
    console.error("Upload error:", error);
  }
}
