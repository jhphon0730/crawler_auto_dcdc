"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { FileIcon, FileTextIcon, ImageIcon, Loader2 } from "lucide-react"

import {
	type Post,
	GetPosts,
} from "@/api/post";

export default function PostsPage() {
  // 상태 관리
  const [posts, setPosts] = useState<Post[]>([])
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [totalPosts, setTotalPosts] = useState(0)
  const [loading, setLoading] = useState(true)

  // 총 페이지 수 계산
  const totalPages = Math.ceil(totalPosts / pageSize)

  // 데이터 로딩
  useEffect(() => {
    const loadPosts = async () => {
      setLoading(true)
      try {
        // 게시글 데이터 가져오기
        const { posts, post_count } = await GetPosts(currentPage, pageSize)
        setPosts(posts)
        setTotalPosts(post_count)
      } catch (error) {
        console.error("데이터 로딩 오류:", error)
      } finally {
        setLoading(false)
      }
    }

    loadPosts()
  }, [currentPage, pageSize])

  // 페이지 변경 핸들러
  const handlePageChange = (page: number) => {
    setCurrentPage(page)
  }

  // 페이지 크기 변경 핸들러
  const handlePageSizeChange = (value: string) => {
    setPageSize(Number(value))
    setCurrentPage(1) // 페이지 크기가 변경되면 첫 페이지로 이동
  }

  // 데이터 타입에 따른 아이콘 표시
  const getDataTypeIcon = (dataType: string) => {
    switch (dataType) {
      case "icon_movie":
        return <FileIcon className="h-4 w-4 text-blue-500" />
      case "icon_pic":
        return <ImageIcon className="h-4 w-4 text-green-500" />
      case "icon_txt":
        return <FileTextIcon className="h-4 w-4 text-gray-500" />
      default:
        return null
    }
  }

  // 데이터 타입에 따른 라벨 표시
  const getDataTypeLabel = (dataType: string) => {
    switch (dataType) {
      case "icon_movie":
        return "동영상"
      case "icon_pic":
        return "이미지"
      case "icon_txt":
        return "텍스트"
      default:
        return "기타"
    }
  }

  // 페이지네이션 링크 생성
  const renderPaginationLinks = () => {
    const links = []

    // 처음 페이지
    if (currentPage > 3) {
      links.push(
        <PaginationItem key="first">
          <PaginationLink onClick={() => handlePageChange(1)}>1</PaginationLink>
        </PaginationItem>,
      )

      if (currentPage > 4) {
        links.push(
          <PaginationItem key="ellipsis-start">
            <PaginationEllipsis />
          </PaginationItem>,
        )
      }
    }

    // 현재 페이지 주변 페이지
    for (let i = Math.max(1, currentPage - 2); i <= Math.min(totalPages, currentPage + 2); i++) {
      links.push(
        <PaginationItem key={i}>
          <PaginationLink isActive={currentPage === i} onClick={() => handlePageChange(i)}>
            {i}
          </PaginationLink>
        </PaginationItem>,
      )
    }

    // 마지막 페이지
    if (currentPage < totalPages - 2) {
      if (currentPage < totalPages - 3) {
        links.push(
          <PaginationItem key="ellipsis-end">
            <PaginationEllipsis />
          </PaginationItem>,
        )
      }

      links.push(
        <PaginationItem key="last">
          <PaginationLink onClick={() => handlePageChange(totalPages)}>{totalPages}</PaginationLink>
        </PaginationItem>,
      )
    }

    return links
  }

  return (
    <div className="container mx-auto py-10">
      <div className="flex justify-between items-center mb-4">
        <div className="text-sm text-muted-foreground">
          {!loading && totalPosts > 0 && (
            <>
              전체 {totalPosts}개 중 {(currentPage - 1) * pageSize + 1}-{Math.min(currentPage * pageSize, totalPosts)}개
              표시 중
            </>
          )}
        </div>
        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground">페이지당 항목:</span>
          <Select value={pageSize.toString()} onValueChange={handlePageSizeChange}>
            <SelectTrigger className="w-[80px]">
              <SelectValue placeholder={pageSize.toString()} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="5">5</SelectItem>
              <SelectItem value="10">10</SelectItem>
              <SelectItem value="20">20</SelectItem>
              <SelectItem value="50">50</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="border rounded-md">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">번호</TableHead>
              <TableHead>제목</TableHead>
              <TableHead className="w-[150px]">작성자</TableHead>
              <TableHead className="w-[180px]">작성일</TableHead>
              <TableHead className="w-[100px]">유형</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={5} className="h-24 text-center">
                  <div className="flex justify-center items-center">
                    <Loader2 className="h-6 w-6 animate-spin mr-2" />
                    데이터를 불러오는 중...
                  </div>
                </TableCell>
              </TableRow>
            ) : posts.length > 0 ? (
              posts.map((post) => (
                <TableRow key={post.post_number} className="cursor-pointer hover:bg-muted/50">
                  <TableCell className="font-medium">{post.post_number}</TableCell>
                  <TableCell>{post.title}</TableCell>
                  <TableCell>{post.writer}</TableCell>
                  <TableCell>{post.write_date}</TableCell>
                  <TableCell>
                    <div className="flex items-center gap-1">
                      {getDataTypeIcon(post.data_type)}
                      <Badge variant="outline">{getDataTypeLabel(post.data_type)}</Badge>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-4">
                  게시글이 없습니다.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {!loading && totalPages > 0 && (
        <div className="mt-4 cursor-pointer">
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => handlePageChange(Math.max(1, currentPage - 1))}
                  className={currentPage === 1 ? "pointer-events-none opacity-50" : ""}
                />
              </PaginationItem>

              {renderPaginationLinks()}

              <PaginationItem>
                <PaginationNext
                  onClick={() => handlePageChange(Math.min(totalPages, currentPage + 1))}
                  className={currentPage === totalPages ? "pointer-events-none opacity-50" : ""}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      )}
    </div>
  )
}

