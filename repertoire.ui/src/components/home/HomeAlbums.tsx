import { Group, ScrollArea, Skeleton, Stack } from '@mantine/core'
import { useGetAlbumsQuery } from '../../state/api/albumsApi.ts'
import HomeAlbumCard from './HomeAlbumCard.tsx'
import { forwardRef } from 'react'

const HomeAlbums = forwardRef<HTMLDivElement, unknown>((_, ref) => {
  const { data: albums, isLoading } = useGetAlbumsQuery({
    pageSize: 20,
    currentPage: 1
  })

  return (
    <ScrollArea ref={ref} scrollbars={'x'} offsetScrollbars={'x'} scrollbarSize={7}>
      <Group wrap={'nowrap'} px={'xl'} align={'start'}>
        {isLoading
          ? Array.from(Array(10)).map((_, i) => (
              <Stack key={i} gap={'xs'} align={'center'}>
                <Skeleton
                  radius={'lg'}
                  h={150}
                  w={150}
                  style={(theme) => ({ boxShadow: theme.shadows.md })}
                />
                <Stack gap={0} align={'center'}>
                  <Skeleton w={100} h={15} mb={4} />
                  <Skeleton w={60} h={10} />
                </Stack>
              </Stack>
            ))
          : albums.models.map((album) => <HomeAlbumCard key={album.id} album={album} />)}
      </Group>
    </ScrollArea>
  )
})

HomeAlbums.displayName = 'HomeAlbums'

export default HomeAlbums
