export type Post = {
	post_number: number;
	title: string;
	content: string;
	writer: string;
	write_date: string;
	data_type: string;
}

type GetPostsResponse = {
	post_count: number;
	posts: Post[];
}

export const GetPosts = async (page = 1, limit = 10): Promise<GetPostsResponse> => {
	try {
		const res = await fetch(`http://localhost:8080/api/posts?page=${page}&limit=${limit}`);

		if (!res.ok) {
			throw new Error('Failed to fetch posts');
		}
		return await res.json();
	} catch (e) {
    return { post_count: 0, posts: [] }
	}
}
