export type Brand = {
  id: number;
  name: string;
};

export type Model = {
  id: number;
  name: string;
};

export type Category = {
  id: number;
  parent?: string;
  name: string;
  path?: string;
};

export type ModelYear = {
  startyear: number;
  endyear: number;
};

export type UserInput = {
  brand?: Brand;
  model?: Model;
  year?: ModelYear;
  category?: Category;
};

export type Motorcycle = {
  id: number;
  brand: string;
  model: string;
  start_year: number;
  end_year: number;
};

export type Product = {
  id: string;
  name: string;
  for_brand: string;
  description: string;
  category_id: number;
  category_path: string;
  is_universal: boolean;
  motorcycles: Motorcycle[];
  importer_name: string;
};
