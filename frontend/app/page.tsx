"use client";

import { DropDown } from "@/components/dropdown";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Result, tryCatch } from "@/utils/trycatch";
import { Brand, Category, Model, ModelYear, Product, UserInput } from "@/utils/types";
import axios from "axios";

import { useEffect, useState } from "react";

export default function Home() {
  const [userInput, setUserInput] = useState<UserInput>({
    brand: undefined,
    model: undefined,
    year: undefined,
    category: undefined,
  });

  const [brands, setBrands] = useState<Brand[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [models, setModels] = useState<Model[]>([]);
  const [years, setYears] = useState<ModelYear[]>([]);
  // Fetching brands
  useEffect(() => {
    async function fetchBrands() {
      const { data, error } = await getBrands();
      if (error !== null) {
        console.error(error);
        return;
      }

      setBrands(data);
    }

    async function fetchCategories() {
      const { data, error } = await getCategories();
      if (error !== null) {
        console.error(error);
        return;
      }

      setCategories(data);
    }

    fetchBrands();
    fetchCategories();
  }, []);

  useEffect(() => {
    async function fetchModels() {
      if (!userInput.brand) return;
      const { data, error } = await getModelsByBrand(userInput.brand);

      if (error !== null) {
        console.error(error);
        return;
      }

      setModels(data);
    }

    fetchModels();

    setUserInput({
      brand: userInput.brand,
    });
  }, [userInput.brand]);

  useEffect(() => {
    async function fetchYears() {
      if (!userInput.brand || !userInput.model) return;
      const { data, error } = await getYears(userInput.brand, userInput.model);
      if (error !== null) {
        console.error(error);
        return;
      }

      setYears(data);
    }

    fetchYears();

    setUserInput({
      brand: userInput.brand,
      model: userInput.model,
    });
  }, [userInput.brand, userInput.model]);

  async function search() {
    const { data, error } = await getFilteredProducts(userInput);
    if (error !== null) {
      console.error("Error when fetching products: ", error);
      return;
    }

    console.log(data);
  }

  return (
    <>
      <div className="flex flex-col items-center justify-center w-full gap-8 py-6">
        <Card className="w-full max-w-6xl">
          <CardHeader>
            <CardTitle>Sök efter delar</CardTitle>
            <CardDescription>Sök och sortera efter delar till både märke samt specifik cross</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-5 grid-rows-1 gap-4 px-8 py-6 rounded-md bg-muted">
              <DropDown
                label="Märke"
                placeholder="Välj märke"
                disabled={false}
                input={userInput}
                setInput={setUserInput}
                type="brand"
                data={brands}
                getOptionLabel={(y) => y.name}
                getOptionValue={(y) => y.name}
              />
              <DropDown
                label="Modell"
                placeholder="Välj modell"
                disabled={userInput?.brand ? false : true}
                input={userInput}
                setInput={setUserInput}
                type="model"
                data={models}
                getOptionLabel={(y) => y.name}
                getOptionValue={(y) => y.name}
              />
              <DropDown
                label="Motorstorlek"
                placeholder="Välj motorstorlek (cc)"
                disabled={userInput?.model ? false : true}
                input={userInput}
                setInput={setUserInput}
                type="year"
                data={years}
                getOptionLabel={(y) => `${y.startyear} - ${y.endyear === 99999 ? "" : y.endyear}`}
                getOptionValue={(y) => `${y.startyear}-${y.endyear}`}
              />
              <DropDown
                label="Kategori"
                placeholder="Välj kategori"
                disabled={false}
                input={userInput}
                setInput={setUserInput}
                type="category"
                data={categories}
                getOptionLabel={(y) => y.name}
                getOptionValue={(y) => y.name}
              />
              <Button className="self-end" variant="default" onClick={search}>
                Sök
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </>
  );
}

/*
  Fetch products
*/
async function getBrands(): Promise<Result<Brand[], Error>> {
  const { data, error } = await tryCatch(axios.get("http://localhost:8000/brands"));
  if (error) {
    console.error("Error when fetching brands: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

async function getModelsByBrand(brand: string): Promise<Result<Model[], Error>> {
  const { data, error } = await tryCatch(axios.get(`http://localhost:8000/brands/${brand}/models`));
  if (error) {
    console.error("Error when fetching models by brand: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

async function getYears(brand: string, model: string): Promise<Result<ModelYear[], Error>> {
  const { data, error } = await tryCatch(axios.get(`http://localhost:8000/brands/${brand}/models/${model}/years`));
  if (error) {
    console.error("Error when fetching years: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

export async function getCategories(): Promise<Result<Category[], Error>> {
  const { data, error } = await tryCatch(axios.get("http://localhost:8000/categories"));
  if (error) {
    console.error("Error when fetching categories: ", error);
    return { data: null, error };
  }
  return { data: data.data, error: null };
}

export async function getFilteredProducts(input: UserInput): Promise<Result<Product[], Error>> {
  const params = new URLSearchParams();

  if (input.brand) params.append("brand", input.brand);
  if (input.model) params.append("model", input.model);
  if (input.year) params.append("year", input.year);
  if (input.category) params.append("category", input.category);

  const { data, error } = await tryCatch(axios.get(`http://localhost:8000/products?${params.toString()}`));

  if (error) {
    console.error("Error fetching filtered products:", error);
    return { data: null, error };
  }

  return { data: data.data, error: null };
}
