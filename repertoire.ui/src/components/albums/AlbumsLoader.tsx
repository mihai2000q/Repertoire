import { Card, Group, Skeleton, Stack } from '@mantine/core'

function AlbumsLoader() {
  return (
    <Group data-testid="albums-loader">
      {Array.from(Array(20)).map((_, i) => (
        <Card key={i} p="sm" shadow="md" h={253} w={175}>
          <Card.Section>
            <Skeleton h={150} />
          </Card.Section>

          <Group justify="space-between" mt="sm" mb="xs">
            <Skeleton w={80} h={12} />
          </Group>

          <Stack gap={4}>
            <Skeleton w={150} h={8} />
            <Skeleton w={150} h={8} />
            <Skeleton w={80} h={8} />
          </Stack>
        </Card>
      ))}
    </Group>
  )
}

export default AlbumsLoader
