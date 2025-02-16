import { ReactElement } from 'react'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { Stack, Text, Title } from '@mantine/core'

function NotFound(): ReactElement {
  useFixedDocumentTitle('Not Found')

  return (
    <Stack px={'xl'}>
      <Title fw={800} order={3}>
        Whoops! Not Found
      </Title>
      <Text>The page you are looking for couldn&#39;t be found</Text>
    </Stack>
  )
}

export default NotFound
