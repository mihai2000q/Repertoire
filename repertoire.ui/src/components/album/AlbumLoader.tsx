import { Card, Divider, Group, Skeleton, Stack } from '@mantine/core'

function AlbumLoader() {
  return (
    <Stack px={'xl'} data-testid={'album-loader'}>
      <Group>
        <Skeleton radius={'lg'} w={'max(12vw, 150px)'} h={'max(12vw, 150px)'} />
        <Stack gap={'xxs'} pt={'10px'}>
          <Skeleton w={75} h={15} />
          <Skeleton w={'max(25vw, 200px)'} h={'max(4vw, 48px)'} my={'xs'} />
          <Group gap={'xxs'}>
            <Skeleton radius={'50%'} w={35} h={35} />
            <Skeleton w={80} h={15} />
            <Skeleton w={35} h={15} />
            <Skeleton w={45} h={15} />
          </Group>
        </Stack>
      </Group>

      <Divider />

      <Card variant={'widget'} mx={'xs'} mb={'lg'} p={0}>
        <Stack px={'md'} pt={'md'} pb={'xs'}>
          <Group>
            <Skeleton w={50} h={15} />
            <Skeleton w={100} h={11} />
          </Group>
          <Stack gap={0}>
            {Array.from(Array(5)).map((_, i) => (
              <Group key={i} px={'md'} py={'xs'}>
                <Skeleton radius={'sm'} w={15} h={16} mr={'xs'} />
                <Skeleton radius={'8px'} w={38} h={38} />
                <Skeleton maw={170} miw={100} h={15} />
              </Group>
            ))}
          </Stack>
        </Stack>
      </Card>
    </Stack>
  )
}

export default AlbumLoader
