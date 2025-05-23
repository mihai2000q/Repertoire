import { Card, Skeleton, Stack } from '@mantine/core'

function SongsLoader() {
  return Array.from(Array(40)).map((_, i) => (
    <Card key={i} p={0} radius={'lg'} shadow={'md'} pb={'md'}>
      <Stack gap={0}>
        <Skeleton radius={'16px'} w={'100%'} h={'unset'} style={{ aspectRatio: 8 / 7 }} />

        <Stack gap={0} px={'sm'} pt={'xs'} pb={6} align={'start'}>
          <Skeleton w={110} h={16} />
          <Skeleton w={80} h={12} my={4} />
        </Stack>
      </Stack>
    </Card>
  ))
}

export default SongsLoader
