import { ReactElement } from 'react'
import useFixedDocumentTitle from "../hooks/useFixedDocumentTitle.ts";

function Unauthorized(): ReactElement {
  useFixedDocumentTitle('Unauthorized')

  return (
    <div>
      <h1>Unauthorized</h1>
      <p>You shouldn&#39;t be here</p>
    </div>
  )
}

export default Unauthorized
