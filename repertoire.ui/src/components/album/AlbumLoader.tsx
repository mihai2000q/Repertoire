import { Group, Skeleton, Stack } from '@mantine/core'

function AlbumLoader() {
  return (
    <Stack>
      <Group align={'center'}>
        <Skeleton radius={'50%'} w={125} h={125} />
        <Skeleton w={200} h={35} />
      </Group>
    </Stack>
  )
}

export default AlbumLoader
