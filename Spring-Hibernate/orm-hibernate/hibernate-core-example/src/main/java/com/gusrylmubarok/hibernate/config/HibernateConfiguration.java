package com.gusrylmubarok.hibernate.config;

import lombok.extern.slf4j.Slf4j;
import org.hibernate.HibernateException;
import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.hibernate.boot.MetadataSources;
import org.hibernate.boot.registry.StandardServiceRegistry;
import org.hibernate.boot.registry.StandardServiceRegistryBuilder;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Slf4j
public class HibernateConfiguration {
    private static final SessionFactory ourSessionFactory;

    static {
        final StandardServiceRegistry registry = new StandardServiceRegistryBuilder()
                .configure() // configura setting from hibernate.cfg.xml
                .build();
        try {
            MetadataSources metadataSources = new MetadataSources(registry);
            // TODO add your mapping here!
            ourSessionFactory = metadataSources.buildMetadata().buildSessionFactory();
        } catch (Throwable ex) {
            StandardServiceRegistryBuilder.destroy(registry);
            throw new ExceptionInInitializerError(ex);
        }
    }

    public static Session getSession() throws HibernateException {
        return ourSessionFactory.openSession();
    }
}
