import { ReactElement } from 'react'
import useFixedDocumentTitle from "../hooks/useFixedDocumentTitle.ts";

function NotFound(): ReactElement {
  useFixedDocumentTitle('Not Found')

  return (
    <div>
      <h1>Whoops! Not Found</h1>
      <p>The page you are looking for couldn&#39;t be found</p>
    </div>
  )
}

export default NotFound
