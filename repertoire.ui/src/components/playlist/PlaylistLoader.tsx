import { Divider, Group, Skeleton, Stack } from '@mantine/core'

function PlaylistLoader() {
  return (
    <Stack data-testid={'playlist-loader'}>
      <Group align={'start'}>
        <Skeleton radius={'lg'} w={150} h={150} />
        <Stack gap={'xxs'} pt={'md'}>
          <Skeleton w={80} h={15} />
          <Skeleton w={200} h={45} my={4} />
          <Skeleton w={85} h={15} />
        </Stack>
      </Group>

      <Divider />

      <Stack p={'sm'} mx={'xs'} mb={'lg'} mt={'xs'}>
        <Group gap={'xs'}>
          <Skeleton w={50} h={15} />
          <Skeleton w={100} h={12} />
        </Group>
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

export default PlaylistLoader
