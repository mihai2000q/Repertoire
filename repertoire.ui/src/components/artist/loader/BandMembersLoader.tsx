import { Card, Group, Skeleton, Stack } from '@mantine/core'

function BandMembersLoader() {
  return (
    <Card variant={'panel'} data-testid={'band-members-loader'} p={0}>
      <Stack gap={0}>
        <Group px={'md'} py={'xs'} gap={'xs'}>
          <Skeleton w={100} h={15} />
        </Group>
        <Group gap={'xs'} wrap={'nowrap'} align={'start'} px={'lg'} pb={'lg'} pt={'xs'}>
          {Array.from(Array(20)).map((_, i) => (
            <Stack key={i} w={75} gap={'xxs'} align={'center'}>
              <Skeleton radius={'50%'} w={55} h={55} />
              <Stack gap={3} align={'center'}>
                <Skeleton w={75} h={10} />
                <Skeleton w={44} h={10} />
              </Stack>
            </Stack>
          ))}
        </Group>
      </Stack>
    </Card>
  )
}

export default BandMembersLoader
