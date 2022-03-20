package com.gusrylmubarok.springbasic.bean;

import com.gusrylmubarok.springbasic.data.Foo;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
public class BeanNameConfiguration {
    @Primary
    @Bean(name = "fooFirst")
    public Foo foo1() {
        return new Foo();
    }

    @Bean(name = "fooSecond")
    public Foo foo2() {
        return new Foo();
    }
}
