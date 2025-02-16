import { Divider, Group, Skeleton, Stack } from '@mantine/core'

function AlbumLoader() {
  return (
    <Stack px={'xl'} data-testid={'album-loader'}>
      <Group align={'start'}>
        <Skeleton radius={'lg'} w={150} h={150} />
        <Stack gap={'xxs'} pt={'10px'}>
          <Skeleton w={75} h={15} />
          <Skeleton w={200} h={45} my={4} />
          <Group gap={'xxs'}>
            <Skeleton radius={'50%'} w={35} h={35} />
            <Skeleton w={80} h={15} />
            <Skeleton w={35} h={15} />
            <Skeleton w={45} h={15} />
          </Group>
        </Stack>
      </Group>

      <Divider />

      <Stack p={'sm'} mx={'xs'} mb={'lg'} mt={'xs'}>
        <Skeleton w={50} h={15} />
        <Stack gap={0}>
          {Array.from(Array(5)).map((_, i) => (
            <Group key={i} px={'md'} py={'xs'}>
              <Skeleton radius={'sm'} w={15} h={16} />
              <Skeleton radius={'8px'} w={38} h={38} />
              <Skeleton maw={170} miw={100} h={15} />
            </Group>
          ))}
        </Stack>
      </Stack>
    </Stack>
  )
}

export default AlbumLoader
