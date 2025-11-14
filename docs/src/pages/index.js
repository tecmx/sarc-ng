import React, { useEffect } from 'react';
import { useHistory } from '@docusaurus/router';

export default function Home() {
  const history = useHistory();

  useEffect(() => {
    // Redirect to introduction page
    history.replace('/content/introduction');
  }, [history]);

  return null;
}
