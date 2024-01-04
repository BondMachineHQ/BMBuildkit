import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';
import ThemedImage from "@theme/ThemedImage";
import useBaseUrl from "@docusaurus/useBaseUrl";
import DocusaurusSvg from '@site/static/img/logo.svg';


function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <header className={clsx('hero hero--primary', styles.heroBanner)}>

            <div className="container">
                <Heading as="h1" className="hero__title">
                    {siteConfig.title}
                </Heading>
                <p className="hero__subtitle">{siteConfig.tagline}</p>
                <div className={styles.buttons}>
                    <DocusaurusSvg className={styles.featureSvg}/>
                </div>
                <div className={styles.buttons}>
                    <Link
                        className="button button--secondary button--lg"
                        to="/docs/intro">
                        Quick-start is just 5min away! ⏱️
                    </Link>
                </div>
            </div>
        </header>
    );
}

export default function Home(): JSX.Element {
    const {siteConfig} = useDocusaurusContext();
    return (
        <Layout
            title={`${siteConfig.title}`}
            description="Evolving FPGA programmining">
            <HomepageHeader/>
            <main>
                <div className="container">
                    <div className="text--center">
                        <Heading as="h2">
                            Built on <a href="https://dagger.io/">dagger.io</a>
                        </Heading>
                    </div>
                </div>
                <HomepageFeatures/>
            </main>
        </Layout>
    );
}
