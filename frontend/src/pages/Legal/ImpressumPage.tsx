import styles from './Legal.module.css';

export const ImpressumPage = () => {
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Legal Notice (Impressum)</h1>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>Information according to § 5 TMG</h2>
        <p className={styles.text}>
          {import.meta.env.VITE_LEGAL_NAME || '[VITE_LEGAL_NAME]'}
          <br />
          {import.meta.env.VITE_LEGAL_ADDRESS || '[VITE_LEGAL_ADDRESS]'},{' '}
          {import.meta.env.VITE_LEGAL_CITY || '[VITE_LEGAL_CITY]'}
          <br />
          {import.meta.env.VITE_LEGAL_COUNTRY || '[VITE_LEGAL_COUNTRY]'}
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>Contact</h2>
        <p className={styles.text}>
          E-Mail: {import.meta.env.VITE_LEGAL_EMAIL || '[VITE_LEGAL_EMAIL]'}
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>Responsible for Content (V.i.S.d.P.)</h2>
        <p className={styles.text}>
          Responsible for content according to § 55 paragraph 2 RStV:
          <br />
          {import.meta.env.VITE_LEGAL_NAME || '[VITE_LEGAL_NAME]'}
          <br />
          {import.meta.env.VITE_LEGAL_ADDRESS || '[VITE_LEGAL_ADDRESS]'},{' '}
          {import.meta.env.VITE_LEGAL_CITY || '[VITE_LEGAL_CITY]'}
          <br />
          {import.meta.env.VITE_LEGAL_COUNTRY || '[VITE_LEGAL_COUNTRY]'}
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>Notice</h2>
        <p className={styles.text}>
          This is a private project for demonstration and educational purposes (Dies ist ein
          privates Projekt zu Demonstrations- und Bildungszwecken). No commercial purpose is pursued
          (Es wird kein wirtschaftlicher Zweck verfolgt).
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>EU Dispute Resolution</h2>
        <p className={styles.text}>
          As this is a private offer and no economic business transaction takes place, we are not
          willing or obliged to participate in dispute resolution proceedings before a consumer
          arbitration board. The EU Online Dispute Resolution (ODR) platform has been discontinued.
        </p>
      </section>
    </div>
  );
};
