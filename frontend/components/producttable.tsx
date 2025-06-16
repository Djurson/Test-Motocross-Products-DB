import { Product } from "@/utils/types";
import { Table, TableBody, TableCell, TableFooter, TableHead, TableHeader, TableRow } from "./ui/table";
import { useState } from "react";
import { Pagination, PaginationContent, PaginationEllipsis, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "./ui/pagination";

const PAGE_SIZE = 30;

export default function ProductTable({ products }: { products: Product[] }) {
  const [currentPage, setCurrentPage] = useState(1);

  const totalPages = Math.ceil(products.length / PAGE_SIZE);

  const startIndex = (currentPage - 1) * PAGE_SIZE;
  const currentProducts = products.slice(startIndex, startIndex + PAGE_SIZE);

  const handlePrev = () => setCurrentPage((p) => Math.max(p - 1, 1));
  const handleNext = () => setCurrentPage((p) => Math.min(p + 1, totalPages));
  const handlePageClick = (page: number) => setCurrentPage(page);

  // För enkelhetens skull visar vi max 5 sidor (kan justeras)
  const pageNumbers = [];
  let startPage = Math.max(1, currentPage - 2);
  let endPage = Math.min(totalPages, currentPage + 2);
  for (let i = startPage; i <= endPage; i++) {
    pageNumbers.push(i);
  }

  return (
    <>
      <Table className="w-full max-w-[90vw]">
        <TableHeader>
          <TableRow>
            <TableHead>Artikelnummer</TableHead>
            <TableHead>Leverantör</TableHead>
            <TableHead>Namn</TableHead>
            <TableHead>Motorcykel märke</TableHead>
            <TableHead>Motorcykel modell</TableHead>
            <TableHead>Under kategori</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {currentProducts.map((product) => (
            <TableRow key={product.id}>
              <TableCell className="font-medium">{product.id}</TableCell>
              <TableCell>{product.importer_name}</TableCell>
              <TableCell>{product.name}</TableCell>
              <TableCell>{product.for_brand}</TableCell>
              <TableCell>
                {product.is_universal ? (
                  <>Universal</>
                ) : (
                  <>
                    {product.motorcycles.length >= 10 ? (
                      <>Fler än 10 st modeller</>
                    ) : (
                      <>
                        {product.motorcycles.map((moto) => {
                          return (
                            <span key={moto.id}>
                              {moto.model} {moto.start_year}-{moto.end_year === 9999 ? "" : moto.end_year}, <br />
                            </span>
                          );
                        })}
                      </>
                    )}
                  </>
                )}
              </TableCell>
              <TableCell>{product.category_path}</TableCell>
            </TableRow>
          ))}
        </TableBody>
        <TableFooter>
          <TableRow>
            <TableCell colSpan={6} className="font-bold text-center">
              Totalt: {products.length} produkter hittade
            </TableCell>
          </TableRow>
        </TableFooter>
      </Table>

      {products.length > 0 && products.length > 30 && (
        <Pagination>
          <PaginationContent>
            {currentPage !== 1 && (
              <PaginationItem>
                <PaginationPrevious onClick={handlePrev} />
              </PaginationItem>
            )}
            {startPage > 1 && (
              <>
                <PaginationItem>
                  <PaginationLink onClick={() => handlePageClick(1)}>1</PaginationLink>
                </PaginationItem>
                {startPage > 2 && (
                  <PaginationItem>
                    <PaginationEllipsis />
                  </PaginationItem>
                )}
              </>
            )}

            {pageNumbers.map((page) => (
              <PaginationItem key={page}>
                <PaginationLink isActive={page === currentPage} onClick={() => handlePageClick(page)}>
                  {page}
                </PaginationLink>
              </PaginationItem>
            ))}

            {endPage < totalPages && (
              <>
                {endPage < totalPages - 1 && (
                  <PaginationItem>
                    <PaginationEllipsis />
                  </PaginationItem>
                )}
                <PaginationItem>
                  <PaginationLink onClick={() => handlePageClick(totalPages)}>{totalPages}</PaginationLink>
                </PaginationItem>
              </>
            )}

            {currentPage !== totalPages && (
              <PaginationItem>
                <PaginationNext onClick={handleNext} />
              </PaginationItem>
            )}
          </PaginationContent>
        </Pagination>
      )}
    </>
  );
}
