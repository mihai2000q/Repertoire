import { ReactElement } from 'react'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { Stack, Text, Title } from '@mantine/core'

function Unauthorized(): ReactElement {
  useFixedDocumentTitle('Unauthorized')

  return (
    <Stack px={'xl'}>
      <Title fw={800} order={3}>
        Unauthorized
      </Title>
      <Text>You shouldn&#39;t be here</Text>
    </Stack>
  )
}

export default Unauthorized
