package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.bean.BeanNameConfiguration;
import com.gusrylmubarok.springbasic.data.Foo;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class BeanNameConfigurationTest {
    private ApplicationContext applicationContext;

    @BeforeEach
    void setUp() {
        applicationContext = new AnnotationConfigApplicationContext(BeanNameConfiguration.class);
    }

    @Test
    void testBeanName() {
        Foo foo = applicationContext.getBean(Foo.class);
        Foo foo1 = applicationContext.getBean("fooFirst", Foo.class);
        Foo foo2 = applicationContext.getBean("fooSecond", Foo.class);

        Assertions.assertSame(foo, foo1);
        Assertions.assertNotSame(foo1, foo2);
    }
}
